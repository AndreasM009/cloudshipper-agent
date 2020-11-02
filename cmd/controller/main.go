package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	natsevents "github.com/andreasM009/cloudshipper-agent/pkg/controller/events/nats"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/configuration"
)

func main() {
	runtime := configuration.NewRuntime()
	runtime.FromFlags()
	// load commands
	azure.LoadAzureCommands()

	// Connection to Nats Server
	natsConnection, err := channel.NewNatsConnection(runtime.NatsServerConnectionStrings, runtime.NatsConnectionName)
	if err != nil {
		log.Panic(err)
	}

	if err := channel.GetNatsConnectionPoolInstance().Add(runtime.NatsConnectionName, natsConnection); err != nil {
		log.Panic(err)
	}

	// connection to Nats Streaming server
	snatsConnection, err := channel.NewNatsStreamingConnectionWithPooledConnection(
		runtime.NatsConnectionName, runtime.NatsStreamingClusterID, runtime.NatsClientID)

	if err != nil {
		log.Panic(err)
	}

	if err := channel.GetNatsStreamingConnectionPoolInstance().Add(runtime.NatsStreamingClusterID, runtime.NatsClientID, snatsConnection); err != nil {
		log.Panic(err)
	}

	// job input stream
	input, err := channel.NewNatsStreamingChannelFromPool(runtime.NatsInputSubscription, runtime.NatsStreamingClusterID, runtime.NatsClientID)

	if err != nil {
		log.Panic(err)
	}

	defer input.Close()

	// publish events channel
	publish, err := channel.NewNatsStreamingChannelFromPool(runtime.NatsPublishSubscription, runtime.NatsStreamingClusterID, runtime.NatsClientID)

	if err != nil {
		log.Panic(err)
	}

	defer publish.Close()

	// publish events forwarder
	forwarder := natsevents.NewNatsStreamEventForwarder(publish, nil)

	ctx, cancel := context.WithCancel(context.Background())
	listener := controller.NewJobListener(input, forwarder)

	err = listener.StartListeningAsync(ctx)
	if err != nil {
		log.Panic(err)
	}

	signalchannel := make(chan os.Signal, 1)
	signal.Notify(signalchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case sig := <-signalchannel:
		log.Println(fmt.Sprintf("Received signal: %d", sig))
		cancel()
		<-listener.Done()
		log.Println("controller ended")
	case <-listener.Done():
		log.Println("listener ended")
	}
}
