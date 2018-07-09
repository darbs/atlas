package main

import (
	"context"
	"math/rand"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/darbs/atlas/internal"
	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-constants/constants"
	"github.com/darbs/barbatos-fwk/database"
	"github.com/darbs/barbatos-fwk/messenger"
	log "github.com/sirupsen/logrus"
)

func init() {
	loglevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(err)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(loglevel)
}

func tearDown(cancel context.CancelFunc) {
	log.Println("Atlas client shutting down")
	internal.StopComm()
	cancel()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	clientId := strconv.Itoa(rand.Int())
	log.Printf("Atlas test client %v", clientId)
	////////////
	// test chunk of code
	////////////

	ctx, cancel := context.WithCancel(context.Background())
	defer tearDown(cancel)

	internal.StartComm(ctx)

	log.Debugf("Client setting up rpc")

	userLocaleChan := make(chan string) // emulate api req
	msgConn := internal.Connection()
	locale := "ABC123DEF123"
	rpc := "rpc-client-" + clientId

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
	log.Debugf("Client requesting rpc")

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

		/////////
		// entity locale update
		log.Debugf("Client listening to locale")
		localeChan, _ := msgConn.Listen(
			constants.AtlasLocaleExchange,
			messenger.ExchangeKindTopic,
			constants.LocaleUpdateKey + "." + locale.Id,
			constants.AtlasLocaleUpdateQueue,
		)

		for {
			msg := <-localeChan
			log.Debugf("Client locale listener rcv: %v", msg)
		}
	}()
	////////////
	////////////
	////////////

	<-ctx.Done()
}