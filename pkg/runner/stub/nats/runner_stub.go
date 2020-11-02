package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreasM009/cloudshipper-agent/pkg/requests"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	base "github.com/andreasM009/cloudshipper-agent/pkg/runner/stub"
	natsio "github.com/nats-io/nats.go"
)

// RunnerStub nats runner stub
type RunnerStub struct {
	channel             *channel.NatsChannel
	healthyHandler      base.HealthyHandler
	cancellationHandler base.CancellationHandler
}

// NewRunnerStub instance
func NewRunnerStub(channel *channel.NatsChannel) *RunnerStub {
	return &RunnerStub{
		channel: channel,
	}
}

// SetHandler callback handlers
func (stub *RunnerStub) SetHandler(healthy base.HealthyHandler, cancellation base.CancellationHandler) error {
	stub.healthyHandler = healthy
	stub.cancellationHandler = cancellation
	return nil
}

// StartListeningAndBlock starts listening for controller messages and block calling thread
func (stub *RunnerStub) StartListeningAndBlock(ctx context.Context, commandrunner <-chan int) error {

	stub.channel.NatsNativeConn.Subscribe(stub.channel.NatsPublishName, func(msg *natsio.Msg) {
		carrier := requests.RequestCarrier{}

		if err := json.Unmarshal(msg.Data, &carrier); err != nil {
			log := fmt.Errorf("[Error]: stub: unable to unmarshal command from json: %s", err)
			fmt.Println(log)
		}

		switch carrier.CarrierForType {
		case requests.RunnerCancel:
			stub.cancellationHandler()
		case requests.RunnerHealth:
			isHealthy := stub.healthyHandler()
			rply := requests.NewRunnerHealthRequest()
			rply.IsHealthy = isHealthy
			carrier.Data = rply

			if json, err := json.Marshal(carrier); err != nil {
				msg.Respond(json)
			} else {
				log := fmt.Errorf("[Error]: stub: unable to marshal command from json: %s", err)
				fmt.Println(log)
			}
		}
	})

	select {
	case <-commandrunner:
	}
	return nil
}
