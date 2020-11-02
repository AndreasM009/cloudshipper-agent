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
	NatsStreamingClusterID      string
	NatsInputSubscription       string
	NatsPublishSubscription     string
	IsKubernetes                bool
	NatsConnectionName          string
	NatsClientID                string
}

var usageStr = `
Usage: runner [options]
Options:
	-m <mode>							Debug or Kubernetes
	-s <url>							NATS server URL(s) (separated by comma)
	-cluster-id <cluster id>			NATS streaming cluster id
	-q <nats job input subscription>	NATS job input subscription
	-publish-subscription 				NATS subscription to publish all events
`
var mode, natsServerURL, clusterID, inputSubscription, publishSubscription string
var theRuntime *Runtime

// NewRuntime new instance
func NewRuntime() *Runtime {
	if nil == theRuntime {
		flag.StringVar(&mode, "m", "Kubernetes", "")
		flag.StringVar(&natsServerURL, "s", "", "")
		flag.StringVar(&clusterID, "cluster-id", "", "")
		flag.StringVar(&inputSubscription, "q", "", "")
		flag.StringVar(&publishSubscription, "publish-subscription", "", "")
		theRuntime = &Runtime{}
	}
	return theRuntime
}

// FromFlags create
func (r *Runtime) FromFlags() {
	flag.Usage = usage
	flag.Parse()

	clientID := fmt.Sprintf("cs-ag-cntrl-%s", uuid.New().String())
	connectionName := clientID

	r.NatsConnectionName = connectionName
	r.NatsClientID = clientID

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

	if clusterID == "" {
		log.Println("no nats streaming cluster id specified")
		flag.Usage()
	}

	if inputSubscription == "" {
		log.Println("no nats job input subscription specified")
		flag.Usage()
	}

	if publishSubscription == "" {
		log.Println("no nats subscription to publish events specified")
		flag.Usage()
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
	r.NatsPublishSubscription = publishSubscription
}

func (r *Runtime) loadDebug() {
	if natsServerURL == "" {
		log.Println("no nats server urls specified")
		flag.Usage()
	}

	if clusterID == "" {
		log.Println("no nats streaming cluster id specified")
		flag.Usage()
	}

	if inputSubscription == "" {
		log.Println("no nats job input subscription specified")
		flag.Usage()
	}

	if publishSubscription == "" {
		log.Println("no nats subscription to publish events specified")
		flag.Usage()
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
	r.NatsPublishSubscription = publishSubscription
}
