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
	userLocaleChan := make(chan string) // emulate api req
	msgConn := internal.Connection()
	locale := "ABC123DEF123"
	rpc := "rpc-123-ABC"

	// unique rpc pubsub
	rpcChan, err := msgConn.Listen(
		constants.AtlasCommandExchange,
		messenger.ExchangeKindDirect,
		rpc,
		rpc,
	)

	if err != nil {
		log.Fatalf("Failed to listen to queue - "+rpc+": %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			msg := <-rpcChan
			log.Debugf("RPC publisher rcv: %v", msg)
			userLocaleChan <- msg.Data
		}
	}()

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

	/////////
	// entity locale update
	localeChan, err := msgConn.Listen(
		constants.AtlasLocaleExchange,
		messenger.ExchangeKindTopic,
		constants.LocaleUpdateKey,
		constants.AtlasLocaleUpdateQueue,
	)

	go func() {
		for {
			msg := <-localeChan
			log.Debugf("Client locale listener rcv: %v", msg)
		}
	}()

	/////////
	// client example listener
	go func() {
		var resp internal.ActionResponse
		raw := <-userLocaleChan
		err := json.Unmarshal([]byte(raw), &resp)
		if err != nil {
			log.Debugf("OOPSIE DAISY: %v", err)
			return
		}

		if resp.Type == internal.ActionStatusError {
			log.Debugf("RESP ERROR: %v", err)
			os.Exit(1)
		}

		var locale model.Locale
		localeData, _ := json.Marshal(resp.Data)
		json.Unmarshal(localeData, &locale)
		log.Debugf("Client received locale %v", locale)

		numOfEntities := 5
		for numOfEntities > 0 {
			ent := model.Entity{
				Id:        database.GetNewObjectId(),
				LocaleId:  locale.Id,
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

			numOfEntities--
		}
	}()
	////////////
	////////////
	////////////

	<-ctx.Done()
}
