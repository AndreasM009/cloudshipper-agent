package configuration

import (
	"flag"
	"log"
	"strings"
)

// Runtime configuration
type Runtime struct {
	NatsServerConnectionStrings []string
	NatsStreamingClusterID      string
	NatsInputSubscription       string
	IsKubernetes                bool
}

var usageStr = `
Usage: runner [options]
Options:
	-m <mode>							Debug or Kubernetes
	-s <url>							NATS server URL(s) (separated by comma)
	-cluster-id <cluster id>			NATS streaming cluster id
	-q <nats job q>						NATS job subscription
`
var mode, natsServerURL, clusterID, inputSubscription string
var theRuntime *Runtime

// NewRuntime new instance
func NewRuntime() *Runtime {
	if nil == theRuntime {
		flag.StringVar(&mode, "m", "Kubernetes", "")
		flag.StringVar(&natsServerURL, "s", "", "")
		flag.StringVar(&clusterID, "cluster-id", "", "")
		flag.StringVar(&inputSubscription, "q", "", "")
		theRuntime = &Runtime{}
	}
	return theRuntime
}

// FromFlags create
func (r *Runtime) FromFlags() {
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

	if clusterID == "" {
		log.Println("no nats streaming cluster id specified")
		flag.Usage()
	}

	if inputSubscription == "" {
		log.Println("no nats job input subscription specified")
		flag.Usage()
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
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

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
}
