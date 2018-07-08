package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/darbs/atlas/internal"
	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-constants/constants"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
	"github.com/darbs/barbatos-fwk/messenger"
	"github.com/sirupsen/logrus"
)

var (
	conf config.Configuration
	log  = logrus.WithFields(logrus.Fields{
		"Component": "Atlas",
	})
)

func tearDown(cancel context.CancelFunc) {
	log.Println("Atlas shutting down")
	internal.StopComm()
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

	ctx, cancel := context.WithCancel(context.Background())
	defer tearDown(cancel)

	internal.StartComm(ctx)
	go internal.ListenForEntityUpdate()
	go internal.ListenForRpc()

	////////////
	// test chunk of code
	////////////
	msgConn := internal.Connection()
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

	/////////
	// test rpc action call
	actionData := map[string]interface{}{}
	actionData["name"] = locale
	actionData["area"] = 30

	action := map[string]interface{}{}
	action["Action"] = internal.OpenLocale
	action["ResponseId"] = rpc
	action["Data"] = actionData

	payload, err := json.Marshal(action)
	if err != nil {
		panic(err)
	}

	msgConn.Publish(
		constants.AtlasCommandExchange,
		messenger.ExchangeKindDirect,
		constants.AtlasRpcKey,
		payload,
	)

	go func() {
		for {
			msg := <-msgChan
			log.Debugf("RPC publisher rcv: %v", msg)
		}
	}()

	// todo publish entity updates
	////////

	for range time.Tick(time.Second * 5) {
		ent := model.Entity{
			Id:        database.GetNewObjectId(),
			LocaleId:  locale,
			Altitude:  4567,
			Longitude: 1234,
			Latitude:  1234,
			Health:    100,
			Mobile:    false,
			Timestamp: time.Now().UTC(),
		}

		payload, _ := json.Marshal(ent)
		msgConn.Publish(
			constants.AtlasEntityExchange,
			messenger.ExchangeKindTopic,
			constants.LocationUpdateKey,
			payload,
		)
	}
	////////////
	////////////
	////////////

	<-ctx.Done()
}
