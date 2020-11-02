package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	"gopkg.in/yaml.v2"
)

var usageStr = `
Usage: runner [options]
Options:
	-s <url>							NATS server URL(s) (separated by comma)
	-q <job queue>						NATS streaming channel for jobs to enqueue
	-cluster-id <cluster id>				NATS streaming server cluster id
	-yaml-definition <yaml definition>	Path to Yaml definition file
	-yaml-parameters					Path to yaml parameters file
`
var natsServerURL, queue, yamldef, parameters, clusterid string

func usage() {
	log.Fatalf(usageStr)
}

type deploymentJob struct {
	TenantID       string            `json:"tenantId"`
	DeploymentName string            `json:"deploymentName"`
	ID             string            `json:"id"`
	DefinitionID   string            `json:"definitionId"`
	Yaml           string            `json:"yaml"`
	Parameters     map[string]string `json:"parameters"`
	LiveStreamName string            `json:"liveStreamName"`
}

func main() {
	flag.StringVar(&natsServerURL, "s", "", "")
	flag.StringVar(&queue, "q", "", "")
	flag.StringVar(&yamldef, "yaml-definition", "", "")
	flag.StringVar(&parameters, "yaml-parameters", "", "")
	flag.StringVar(&clusterid, "cluster-id", "", "")

	flag.Parse()

	channel, err := channel.NewNatsStreamingChannel(queue, strings.Split(natsServerURL, ","), "jobclient", clusterid, "jobclient")
	if err != nil {
		log.Panic(err)
	}

	defer channel.Close()

	deploymentID := uuid.New().String()

	yamlcontent, err := ioutil.ReadFile(yamldef)
	if err != nil {
		log.Panic(err)
	}

	parametercontent, err := ioutil.ReadFile(parameters)
	if err != nil {
		log.Panic(err)
	}

	parammap := map[string]string{}
	err = yaml.Unmarshal(parametercontent, &parammap)
	if err != nil {
		log.Panic(err)
	}

	job := deploymentJob{
		DeploymentName: fmt.Sprintf("JobClientDeployment-%s", deploymentID),
		TenantID:       "t1",
		ID:             deploymentID,
		DefinitionID:   "1",
		Yaml:           string(yamlcontent),
		LiveStreamName: "agentlive",
		Parameters:     parammap,
	}

	json, err := json.Marshal(job)
	if err != nil {
		log.Panic(err)
	}

	err = channel.SnatNativeConnection.Publish(channel.NatsPublishName, json)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Deployment Job queued!")
}
