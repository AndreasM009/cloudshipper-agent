package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	input, err := channel.NewNatsChannel(runtime.NatsInputSubscription, runtime.NatsServerConnectionStrings, "cs-agent-controller")
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	listener := controller.NewJobListener(input)

	err = listener.StartListeningAsync(ctx)
	if err != nil {
		log.Panic(err)
	}

	signalchannel := make(chan os.Signal, 1)
	signal.Notify(signalchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-signalchannel:
		cancel()
		<-listener.Done()
	case <-listener.Done():
		log.Println("listener ended")
	}
}
