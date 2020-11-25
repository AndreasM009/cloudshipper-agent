package configuration

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	RunnerExecutableFilePath    string
	DebugRunnerChannelName      string
	IsKubernetes                bool
	IsDebug                     bool
	IsStandalone                bool
	NatsConnectionName          string
	NatsClientID                string
	NatsToken                   string
}

var usageStr = `
Usage: runner [options]
Options:
	-m <mode>									Debug, Standalone or Kubernetes
	-s <url>									NATS server URL(s) (separated by comma)
	-cluster-id <cluster id>					NATS streaming cluster id
	-q <nats job input subscription>			NATS job input subscription
	-publish-subscription 						NATS subscription to publish all events
	-rp-standalone Path to runner executable	Path to runner executable for mode Standalone only
	-rcn-debug Name of NATS channel to runner	Name of NATS channel to runner for mode Debug
	-nats-token-filepath						Path to file containing NATS Auth token, if NATS requires authentication
`
var mode, natsServerURL, clusterID, inputSubscription, publishSubscription, runnerexecutable, debugrunnerchannel, natstokenfile string
var theRuntime *Runtime

// GetRuntime new instance
func GetRuntime() *Runtime {
	if nil == theRuntime {
		flag.StringVar(&mode, "m", "Kubernetes", "")
		flag.StringVar(&natsServerURL, "s", "", "")
		flag.StringVar(&clusterID, "cluster-id", "", "")
		flag.StringVar(&inputSubscription, "q", "", "")
		flag.StringVar(&publishSubscription, "publish-subscription", "", "")
		flag.StringVar(&runnerexecutable, "rp-standalone", "", "")
		flag.StringVar(&debugrunnerchannel, "rcn-debug", "", "")
		flag.StringVar(&natstokenfile, "nats-token-filepath", "", "")
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
	} else if strings.ToLower(mode) == "standalone" {
		r.loadStandalone()
	} else if strings.ToLower(mode) == "debug" {
		r.loadDebug()
	} else {
		flag.Usage()
	}

	if natstokenfile != "" {
		data, err := ioutil.ReadFile(natstokenfile)
		if err != nil {
			log.Println(fmt.Sprintf("nats token file specified, but could not be read: %s", err))
			flag.Usage()
		}

		r.NatsToken = string(data)
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
	r.IsDebug = false
	r.IsKubernetes = true
	r.IsStandalone = false
}

func (r *Runtime) loadStandalone() {
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

	if runnerexecutable == "" {
		log.Println("no file path to runner executable specified")
		flag.Usage()
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
	r.NatsPublishSubscription = publishSubscription
	r.RunnerExecutableFilePath = runnerexecutable
	r.IsDebug = false
	r.IsKubernetes = false
	r.IsStandalone = true
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

	if debugrunnerchannel == "" {
		log.Println("no runner channel name specified for debug mode")
	}

	r.NatsServerConnectionStrings = strings.Split(natsServerURL, ",")
	r.NatsStreamingClusterID = clusterID
	r.NatsInputSubscription = inputSubscription
	r.NatsPublishSubscription = publishSubscription
	r.DebugRunnerChannelName = debugrunnerchannel
	r.IsDebug = true
	r.IsKubernetes = false
	r.IsStandalone = false
}
