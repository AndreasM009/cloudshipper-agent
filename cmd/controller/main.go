package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/runner"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/configuration"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/processor"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/definition"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"

	natsforwarder "github.com/andreasM009/cloudshipper-agent/pkg/controller/events/nats"
)

const (
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
      - command: AgentDownloadArtifacts
        with:
          url: https://anmockartifacts.blob.core.windows.net/release/main.zip
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
	runtime := configuration.NewRuntime()
	runtime.FromFlags()

	// load commands
	azure.LoadAzureCommands()

	servicePrincipalName := os.Getenv(envServicePrincipalName)
	servicePrincipalSecret := os.Getenv(envServicePrincipalSecret)
	tenant := os.Getenv(envTenant)
	subscription := os.Getenv(envSubscription)

	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envServicePrincipalName, servicePrincipalName)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envServicePrincipalSecret, servicePrincipalSecret)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envSubscription, subscription)
	yamlTestDeploymentDefinition = strings.ReplaceAll(yamlTestDeploymentDefinition, envTenant, tenant)

	runnerChannel, err := channel.NewNatsChannel(runtime.NatsRunnerChannelName, runtime.NatsServerConnectionStrings, runtime.NatsConnectionName)
	if err != nil {
		log.Panic(err)
	}

	defer runnerChannel.Close()

	liveStream, err := channel.NewNatsStreamingChannel(runtime.NatsLiveStreamName, runtime.NatsServerConnectionStrings, runtime.NatsConnectionName, runtime.NatsStreamingClusterID, runtime.NatsClientID)
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

	forwarder := natsforwarder.NewNatsStreamEventForwarder(liveStream)

	proc := processor.NewCommandProcessor(runnerChannel, runtimeDep, forwarder)

	err = proc.ProcessAsync()
	if err != nil {
		log.Panic(err)
	}

	// start the RunnerPod if runner is hosted in Kubernetes
	k8sclient := runner.CreateKubeClient()

	if k8sclient == nil {
		log.Panic()
	}

	if runtime.IsKubernetes {
		runnerCtx, cancelRunner := context.WithCancel(context.Background())

		k8srunner, err := runner.CreateAndWatchRunnerPod(runnerCtx, k8sclient, runtime.NatsRunnerChannelName)
		if err != nil {
			log.Panic(err)
		}

		select {
		case exitcode := <-proc.Done():
			cancelRunner()
			<-k8srunner.Error()
			fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
			return
		case err := <-k8srunner.Error():
			if err != nil {
				log.Panic(err)
			}
		}
	} else {
		exitcode := <-proc.Done()
		fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
	}
}
