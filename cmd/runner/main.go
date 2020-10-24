package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/settings"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"

	natsproxy "github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy/nats"
	natsstub "github.com/andreasM009/cloudshipper-agent/pkg/runner/stub/nats"

	"github.com/andreasM009/cloudshipper-agent/pkg/runner/commands"
)

const ()

func main() {
	// load commands
	azure.LoadAzureCommands()

	// context with cancel to be able to gracefully shutdown runner
	ctx, cancel := context.WithCancel(context.Background())
	// signals
	sigchannel := make(chan os.Signal)
	signal.Notify(sigchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// nats channel
	natsChannel, err := channel.NewNatsChannel(settings.GetControllerRunnerChannelName(), []string{"localhost:4222"}, fmt.Sprintf("%s-runner", settings.GetControllerRunnerChannelName()))
	if err != nil {
		log.Panic(err)
	}

	defer natsChannel.Close()

	proxy, err := natsproxy.NewNatsControllerProxy(natsChannel)
	if err != nil {
		log.Panic(err)
	}
	defer natsChannel.Close()

	healthyhandler := func() bool { return true }
	cancelHandler := func() bool {
		cancel()
		return true
	}

	stub := natsstub.NewRunnerStub(natsChannel)
	stub.SetHandler(healthyhandler, cancelHandler)

	commandrunner := make(chan int, 1)

	go func() {
		// command result
		result := 0

		// command loop
		for result == 0 {
			cmdStream, err := proxy.GetAgentCommand()
			if err != nil {
				fmt.Println(err)
				commandrunner <- 1
			}

			if nil == cmdStream || len(cmdStream) == 0 {
				fmt.Println("No further command to execute, runner is shutting down")
				commandrunner <- 0
				return
			}

			cmd, err := commands.CreateCommand(cmdStream, proxy)
			if err != nil {
				fmt.Println(err)
				commandrunner <- 1
			}

			if err := cmd.ExecuteAsync(ctx); err != nil {
				fmt.Println(err)
				commandrunner <- 1
			}

			select {
			// cancel?
			case <-ctx.Done():
				exitcode := <-cmd.Done()
				proxy.Report(logs.NewInfoLog(fmt.Sprintf("Canceled by user -> Exitcode: %d", exitcode)))
				commandrunner <- exitcode
				return
			// cancel by system
			case sig := <-sigchannel:
				cancel()
				exitcode := <-cmd.Done()
				proxy.Report(logs.NewInfoLog(fmt.Sprintf("Runner stopped due to signal %s -> Exitcode: %d", sig, exitcode)))
				commandrunner <- exitcode
				return
			// command finished
			case exitcode := <-cmd.Done():
				proxy.Report(logs.NewInfoLog(fmt.Sprintf("Command finished with exitcode %d", exitcode)))

				// error during command execution?
				if 0 != exitcode {
					proxy.ReportError(exitcode)
					commandrunner <- exitcode
					return
				}
			}
		}
	}()

	// Start listening for controller requests and block until all commands are executed
	// or until runner has to cancel
	stub.StartListeningAndBlock(ctx, commandrunner)
}