package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"

	"github.com/andreasM009/nats-library/channel"
)

var usageStr = `
Usage: runner [options]
Options:
	-s <url>							NATS server URL(s) (separated by comma)
	-l <live streaming channel name>	NATS streaming channel for live logs
	-cluster-id <cluster id>			NATS streaming cluster id
	-t <NAT auth token>					NATS auth token if needed
`
var natsServerURL, liveChannelName, clusterID, natsToken string

func usage() {
	log.Fatalf(usageStr)
}

func main() {
	flag.StringVar(&natsServerURL, "s", "", "")
	flag.StringVar(&liveChannelName, "l", "", "")
	flag.StringVar(&clusterID, "cluster-id", "", "")
	flag.StringVar(&natsToken, "t", "", "")

	flag.Usage = usage
	flag.Parse()

	streamingChannel, err := channel.NewNatsStreamingChannel(liveChannelName, strings.Split(natsServerURL, ","), "livestream-client", clusterID, "livestream-client", natsToken)
	if err != nil {
		log.Panic(err)
	}

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	handler := func(msg *stan.Msg) {
		//log.Printf("[#%d] Received: %s\n", msg.Sequence, string(msg.Data))
		data := map[string]interface{}{}

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Println(err)
		}

		if val, ok := data["eventName"]; ok {
			evtName := val.(string)

			if strings.ToLower(evtName) == "deploymentevent" {
				dpl := struct {
					TenantID       string    `json:"tenantId"`
					DeploymentName string    `json:"deploymentName"`
					Started        bool      `json:"started"`
					Finished       bool      `json:"finished"`
					Exitcode       int       `json:"exitcode"`
					Timestamp      time.Time `json:"timestamp"`
				}{}

				if err := json.Unmarshal(msg.Data, &dpl); err == nil {
					if dpl.Started {
						log.Println("*******************************************")
						log.Println("Deployment:", dpl.DeploymentName, "for tenant:", dpl.TenantID, "started at ", dpl.Timestamp)
						log.Println("*******************************************")
					} else if dpl.Finished {
						log.Println("Deployment:", dpl.DeploymentName, "for tenant:", dpl.TenantID, "finished: exitcode ", dpl.Exitcode, " time ", dpl.Timestamp)
					}
				}

			} else if strings.ToLower(evtName) == "commandevent" {
				cmd := struct {
					CommandDisplayName string    `json:"commandDisplayName"`
					Timestamp          time.Time `json:"timestamp"`
					Logs               []struct {
						Message string `json:"message"`
					} `json:"logs"`
				}{}

				if err := json.Unmarshal(msg.Data, &cmd); err == nil {
					fmt.Println(cmd.Timestamp, " ", cmd.CommandDisplayName, ":", cmd.Logs[0].Message)
				}
			}
		}
	}

	sub, err := streamingChannel.SnatNativeConnection.Subscribe(streamingChannel.NatsPublishName, handler, startOpt)
	if err != nil {
		log.Panic(err)
	}

	signalchannel := make(chan os.Signal, 1)
	signal.Notify(signalchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	defer sub.Unsubscribe()
	defer streamingChannel.Close()
	<-signalchannel
}
