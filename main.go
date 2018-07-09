package main

import (
	"context"
	"os/signal"
	"sync"
	"os"

	"github.com/darbs/atlas/internal"
	"github.com/darbs/barbatos-fwk/config"
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
	internal.ActionShutdown()
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

	// todo this is ugly
	// sig catch
	var end_waiter sync.WaitGroup
	var signal_channel chan os.Signal
	end_waiter.Add(1)
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
	////////

	log.Debugf("Shutdown")
	cancel()
}
