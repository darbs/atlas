package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-constants/constants"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/messenger"
	"github.com/globalsign/mgo/bson"
)

var (
	conf config.Configuration
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
		log.Fatalf("Failed to listen to queue - " + constants.AtlasEntityUpdateQueue + ": %v", err)
		os.Exit(1)
	}

	for {
		msg := <-msgChan
		entity, err := model.EntityFromJson(msg.Data)
		if err != nil {
			log.Printf("Error parsing: %v msg: %v/n", err, msg)
			// todo some sort of error tracking here
			continue
		}

		err = entity.Save()
		if err != nil {
			log.Printf("Failed to save entity: %v", err)
			continue
		}

		log.Printf("entity recieved: %v", entity.Altitude)
	}
}

func listenForLocaleUpdate(conn messenger.Connection) {
	msgChan, err := conn.Listen(
		constants.AtlasLocaleExchange,
		messenger.ExchangeKindTopic,
		constants.LocaleUpdateKey,
		constants.AtlasLocaleUpdateQueue,
	)

	if err != nil {
		log.Fatalf("Failed to listen to queue - " + constants.AtlasLocaleUpdateQueue + ": %v", err)
		os.Exit(1)
	}

	for {
		msg := <-msgChan
		log.Printf("raw message: %v", msg)

		entity, err := model.EntityFromJson(msg.Data)
		if err != nil {
			log.Printf("Error parsing: %v msg: %v/n", err, msg)
			// todo some sort of error tracking here
			continue
		}

		// TODO pull out correlationId and replyTo to attach to response message after fetching all entities from current local that are not you
		log.Printf("local request recieved: %v", entity.Altitude)
	}
}

func publishLocaleUpdate(conn messenger.Connection) {

}

func tearDown(cancel context.CancelFunc, connection messenger.Connection) {
	log.Println("Atlas shutting down")
	connection.Stop()
	cancel()
}

func main() {
	log.Println("Initializing Atlas")

	conf := config.GetConfig()
	msgConn := initializeMqConnection(conf.MqEndpoint)

	ctx, cancel := context.WithCancel(context.Background())
	go msgConn.Start(ctx)

	defer tearDown(cancel, msgConn)

	go listenForEntityUpdate(msgConn)

	go listenForLocaleUpdate(msgConn)

	////////////
	// test chunk of code
	////////////
	locale := "ABC123DEF123"

	for range time.Tick(time.Second * 5) {
		ent := model.Entity{
			Id: bson.NewObjectId().String(),
			Locale:    locale,
			Altitude:  4567,
			Longitude: 1234,
			Latitude:  1234,
			Health:    100,
			Mobile:    false,
		}

		payload, _ := json.Marshal(ent)
		//entity based
		msgConn.Publish(
			constants.AtlasEntityExchange,
			messenger.ExchangeKindTopic,
			constants.LocationUpdateKey,
			payload,
		)

		ent2 := model.Entity{
			Id: bson.NewObjectId().String(),
			Locale:    locale,
			Altitude:  9000,
			Longitude: 1234,
			Latitude:  1234,
			Health:    100,
			Mobile:    false,
		}

		payload2, _ := json.Marshal(ent2)
		msgConn.Publish(
			constants.AtlasLocaleExchange,
			messenger.ExchangeKindTopic,
			constants.LocaleUpdateKey,
			payload2,
		)

	}
	////////////
	////////////
	////////////

	<-ctx.Done()
}
