package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/darbs/atlas/internal"
	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-constants/constants"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/messenger"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

var (
	conf config.Configuration
	log = logrus.WithFields(logrus.Fields{
		"Component": "Atlas",
	})
)


func initializeMqConnection(endpoint string) messenger.Connection {
	log.Println("Initiliazing message connection")

	var conf = messenger.Config{
		Url:       endpoint,
		Durable:   true,
		Attempts:  5,
		Delay:     time.Second * 2,
		Threshold: 4,
	}
	var msgConn, err = messenger.GetConnection(conf)
	if err != nil {
		panic(fmt.Errorf("Failed to connect to message queue: %v", err))
	}

	return msgConn
}

func listenForEntityUpdate(conn messenger.Connection) {
	msgChan, err := conn.Listen(
		constants.AtlasEntityExchange,
		messenger.ExchangeKindTopic,
		constants.LocationUpdateKey,
		constants.AtlasEntityUpdateQueue,
	)

	if err != nil {
		log.Fatalf("Failed to listen to queue - "+constants.AtlasEntityUpdateQueue+": %v", err)
		os.Exit(1)
	}

	for {
		msg := <-msgChan
		entity, err := model.EntityFromJson(msg.Data)
		if err != nil {
			log.Info("Error parsing: %v msg: %v/n", err, msg)
			// todo some sort of error tracking here
			continue
		}

		err = entity.Save()
		if err != nil {
			log.Info("Failed to save entity: %v", err)
			continue
		}

		log.Debugf("Entity recieved: %v", entity.Altitude)
	}
}

func listenForRpc(conn messenger.Connection) {
	msgChan, err := conn.Listen(
		constants.AtlasCommandExchange,
		messenger.ExchangeKindDirect,
		constants.AtlasRpcKey,
		constants.AtlasCommandQueue,
	)

	if err != nil {
		log.Fatalf("Failed to listen to queue - "+constants.AtlasEntityUpdateQueue+": %v", err)
		os.Exit(1)
	}

	for {
		var msgrcv messenger.RpcMessage
		msg := <-msgChan
		err := json.Unmarshal([]byte(msg.Data), &msgrcv)
		if err != nil {
			log.Printf("RPC Error: %v", err)
			continue
		}

		log.Debugf("RPC recieved: %v", msg.Data)
		if msgrcv.ResponseId == "" {
			log.Warnf("No rpc response queue: %v", msgrcv)
			continue
		}

		if msgrcv.Action == "" {
			log.Warnf("No provided action: %v", msgrcv)
			continue
		}

		resp := actions.Handler(msgrcv.Action, msgrcv.Data)
		payload, err := json.Marshal(resp)
		if err != nil {
			// todo error response
		}

		err = conn.Publish(
			constants.AtlasCommandExchange,
			messenger.ExchangeKindDirect,
			msgrcv.ResponseId,
			payload,
		)
		if err != nil {
			log.Infof("Error responding to RPC: %v", err)
		}
	}
}

func tearDown(cancel context.CancelFunc, connection messenger.Connection) {
	log.Println("Atlas shutting down")
	connection.Stop()
	cancel()
}

func init() {
	loglevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(loglevel)
}

func main() {
	log.Println("Initializing")

	conf := config.GetConfig()
	msgConn := initializeMqConnection(conf.MqEndpoint)

	ctx, cancel := context.WithCancel(context.Background())
	defer tearDown(cancel, msgConn)

	go msgConn.Start(ctx)

	go listenForEntityUpdate(msgConn)

	go listenForRpc(msgConn)

	//go publishLocaleUpdate(msgConn)

	////////////
	// test chunk of code
	////////////
	locale := "ABC123DEF123"
	rpc := "rpc-123-ABC"

	// unique rpc pubsub
	msgChan, err := msgConn.Listen(
		constants.AtlasCommandExchange,
		messenger.ExchangeKindDirect,
		rpc,
		rpc,
	)

	if err != nil {
		log.Fatalf("Failed to listen to queue - "+constants.AtlasEntityUpdateQueue+": %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			msg := <-msgChan
			log.Debugf("RPC publisher rcv: %v", msg)
		}
	}()


	for range time.Tick(time.Second * 5) {
		ent := model.Entity{
			Id:        bson.NewObjectId().String(),
			Locale:    locale,
			Altitude:  4567,
			Longitude: 1234,
			Latitude:  1234,
			Health:    100,
			Mobile:    false,
		}

		payload, _ := json.Marshal(ent)
		msgConn.Publish(
			constants.AtlasEntityExchange,
			messenger.ExchangeKindTopic,
			constants.LocationUpdateKey,
			payload,
		)

		var mqMsg interface{}
		msg := []byte(`{"Action":"INITIALIZE_LOCALE", "ResponseId": "` + rpc + `", "Data": { "name": "Hello", "Area": 30 }}`)
		err := json.Unmarshal(msg, &mqMsg)
		if err != nil {
			log.Errorf("ERR: %v", err)
		}

		//payload2, _ := json.Marshal(mqMsg)
		msgConn.Publish(
			constants.AtlasCommandExchange,
			messenger.ExchangeKindDirect,
			constants.AtlasRpcKey,
			msg,
		)
	}
	////////////
	////////////
	////////////

	<-ctx.Done()
}
