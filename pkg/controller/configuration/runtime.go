package configuration

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

// Runtime configuration
type Runtime struct {
	NatsServerConnectionStrings []string
	NatsRunnerChannelName       string
	NatsLiveStreamName          string
	NatsStreamingClusterID      string
	NatsConnectionName          string
	NatsClientID                string
	IsKubernetes                bool
}

var usageStr = `
Usage: runner [options]
Options:
	-m <mode>							Debug or Kubernetes
	-s <url>							NATS server URL(s) (separated by comma)
	-c <channel name> 					NATS channel name to controller if mode == Debug
	-l <live streaming channel name>	NATS streaming channel for live logs if mode == Debug
	-cluster-id <cluster id>			NATS streaming cluster id
`
var mode, natsServerURL, runnerChannelName, liveChannelName, clusterID string

// NewRuntime new instance
func NewRuntime() *Runtime {
	return &Runtime{}
}

// FromFlags create
func (r *Runtime) FromFlags() {
	flag.StringVar(&mode, "m", "Kubernetes", "")
	flag.StringVar(&natsServerURL, "s", "", "")
	flag.StringVar(&runnerChannelName, "c", "", "")
	flag.StringVar(&liveChannelName, "l", "", "")
	flag.StringVar(&clusterID, "cluster-id", "", "")

	flag.Usage = usage
	flag.Parse()

	if strings.ToLower(mode) == "kubernetes" {
		r.loadKubernetes()
	} else if strings.ToLower(mode) == "debug" {
		r.loadDebug()
	} else {
		flag.Usage()
	}
}

func usage() {
	log.Fatalf(usageStr)
}

func (r *Runtime) loadKubernetes() {
	r.IsKubernetes = true
	if natsServerURL == "" {
		log.Println("no nats server urls specified")
		flag.Usage()
	}

	if liveChannelName == "" {
		log.Println("no live stream channel name specified")
		flag.Usage()
	}

	if clusterID == "" {
		log.Println("no nats streaming cluster id specified")
		flag.Usage()
	}

	// Guid
	id := uuid.New()

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsRunnerChannelName = id.String()
	r.NatsLiveStreamName = liveChannelName
	r.NatsStreamingClusterID = clusterID
	r.NatsConnectionName = fmt.Sprintf("controller-%s", id.String())
	r.NatsClientID = fmt.Sprintf("controller-%s", id.String())
}

func (r *Runtime) loadDebug() {
	if natsServerURL == "" {
		log.Println("no nats server urls specified")
		flag.Usage()
	}

	if runnerChannelName == "" {
		log.Println("no nats channel name for runner specified")
		flag.Usage()
	}

	if liveChannelName == "" {
		log.Println("no live stream channel name specified")
		flag.Usage()
	}

	if clusterID == "" {
		log.Println("no nats streaming cluster id specified")
		flag.Usage()
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsRunnerChannelName = runnerChannelName
	r.NatsLiveStreamName = liveChannelName
	r.NatsStreamingClusterID = clusterID
	r.NatsConnectionName = "controller-debug"
	r.NatsClientID = "controller-debug"
}
