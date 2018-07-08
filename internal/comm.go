package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-constants/constants"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/messenger"
	log "github.com/sirupsen/logrus"
)

var (
	connection messenger.Connection
	started bool
)

func init() {
	conf := config.GetConfig()
	//msgConn := InitializeConnection(conf.MqEndpoint)
	log.Println("Initializing message connection")

	var err error
	connection, err = messenger.GetConnection(messenger.Config{
		Url:       conf.MqEndpoint,
		Durable:   true,
		Attempts:  5,
		Delay:     time.Second * 2,
		Threshold: 4,
	})
	if err != nil {
		panic(fmt.Errorf("Failed to connect to message queue: %v", err))
	}
}

func StartComm(ctx context.Context) {
	if started != true {
		go connection.Start(ctx)
	}

	started = true
}

func StopComm() {
	connection.Stop()
	started = false
}

func Connection () messenger.Connection{
	return connection
}

func ListenForEntityUpdate() {
	msgChan, err := connection.Listen(
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

		log.Debugf("Entity received: %v", entity)
	}
}

func ListenForRpc() {
	msgChan, err := connection.Listen(
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

		log.Debugf("RPC received: %v", msg.Data)
		if msgrcv.ResponseId == "" {
			log.Warnf("No rpc response queue: %v", msgrcv)
			continue
		}

		if msgrcv.Action == "" {
			log.Warnf("No provided action: %v", msgrcv)
			continue
		}

		resp := ActionHandler(msgrcv.Action, msgrcv.Data)
		payload, err := json.Marshal(resp)
		if err != nil {
			// todo error response
		}

		err = connection.Publish(
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

func BroadcastToLocale (payload map[string] interface{}) error {
	jsonResp, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling response: %v", err)
	}

	err = connection.Publish(
		constants.AtlasLocaleExchange,
		messenger.ExchangeKindTopic,
		constants.LocaleUpdateKey,
		jsonResp,
	)
	return err
}
