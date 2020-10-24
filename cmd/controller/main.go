package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/processor"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/definition"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"

	"github.com/andreasM009/cloudshipper-agent/pkg/runner/settings"

	natsforwarder "github.com/andreasM009/cloudshipper-agent/pkg/controller/events/nats"
)

const (
	envScriptToRun            = "SCRIPT_TO_RUN"
	envServicePrincipalName   = "SERVICEPRINCIPAL_NAME"
	envServicePrincipalSecret = "SERVICEPRINCIPAL_SECRET"
	envTenant                 = "TENANT"
	envSubscription           = "SUBSCRIPTION"
)

var yamlTestDeploymentDefinition = `
jobs:
  testing:
    displayname: test deployment
    working-directory: ./main
    steps:
      - command: AzPowerShellCore
        working-directory: ./main
        displayname: deploy 1st ARM Template
        with:
          arguments: -ResourceGroupName CLOUD-SHIPPER-RG -StorageAccountName anmocktst
          scriptToRun: ./deploy-arm-template.ps1
          subscription: SUBSCRIPTION
          tenant: TENANT
          serviceprincipal: SERVICEPRINCIPAL_NAME
          secret: SERVICEPRINCIPAL_SECRET
      - command: AzPowerShellCore
        working-directory: ./main
        displayname: deploy 2nd ARM Template
        with:
          arguments: -ResourceGroupName CLOUD-SHIPPER-RG -StorageAccountName anmockshp
          scriptToRun: ./deploy-arm-template.ps1
          subscription: SUBSCRIPTION
          tenant: TENANT
          serviceprincipal: SERVICEPRINCIPAL_NAME
          secret: SERVICEPRINCIPAL_SECRET
parameters: {}
`

func main() {
	fmt.Println("Hello controller")

	// load commands
	azure.LoadAzureCommands()

	servicePrincipalName := os.Getenv(envServicePrincipalName)
	servicePrincipalSecret := os.Getenv(envServicePrincipalSecret)
	tenant := os.Getenv(envTenant)
	subscription := os.Getenv(envSubscription)
	natsChannelName := settings.GetControllerRunnerChannelName()
	natsStreamingLiveChannelName := os.Getenv("NATS_LIVE_STREAMING_CHANNEL_NAME")

	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envServicePrincipalName, servicePrincipalName)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envServicePrincipalSecret, servicePrincipalSecret)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envSubscription, subscription)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envTenant, tenant)

	runnerChannel, err := channel.NewNatsChannel(natsChannelName, []string{"localhost:4222"}, "debugging-controller")
	if err != nil {
		log.Panic(err)
	}

	defer runnerChannel.Close()

	liveStream, err := channel.NewNatsStreamingChannel(natsStreamingLiveChannelName, []string{"localhost:4222"}, "debugging-live-controller", "test-cluster", "livestream")
	if err != nil {
		log.Panic(err)
	}

	defer liveStream.Close()

	def, err := definition.NewFromYaml([]byte(yamlTestDeploymentDefinition))
	if err != nil {
		log.Panic(err)
	}

	runtimeDep, err := definition.NewFromDefinition(def, "23", "test", "1", make(map[string]string))
	if err != nil {
		log.Panic(err)
	}

	//forwarder := fakeforwarder.NewFakeEventForwarder()
	forwarder := natsforwarder.NewNatsStreamEventForwarder(liveStream)

	proc := processor.NewCommandProcessor(runnerChannel, runtimeDep, forwarder)

	err = proc.ProcessAsync()
	if err != nil {
		log.Panic(err)
	}

	exitcode := <-proc.Done()
	fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
}
