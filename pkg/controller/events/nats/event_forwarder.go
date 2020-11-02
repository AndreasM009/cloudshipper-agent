package nats

import (
	"encoding/json"
	"fmt"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	base "github.com/andreasM009/cloudshipper-agent/pkg/controller/events"
)

// EventForwarder forwards events to nats stream
type EventForwarder struct {
	stream        *channel.NatsStreamingChannel
	nextForwarder base.EventForwarder
}

// NewNatsStreamEventForwarder new instance
func NewNatsStreamEventForwarder(stream *channel.NatsStreamingChannel, next base.EventForwarder) base.EventForwarder {
	return &EventForwarder{
		stream:        stream,
		nextForwarder: next,
	}
}

// ForwardDeploymentEvent event
func (forwarder *EventForwarder) ForwardDeploymentEvent(evt *base.DeploymentEvent) {
	json, err := json.Marshal(evt)
	if err != nil {
		fmt.Println("Error forwarding event to nats streaming server")
	}

	//fmt.Println(string(json))
	forwarder.stream.SnatNativeConnection.Publish(forwarder.stream.NatsPublishName, json)

	if forwarder.nextForwarder != nil {
		forwarder.nextForwarder.ForwardDeploymentEvent(evt)
	}
}

// ForwardCommandEvent event
func (forwarder *EventForwarder) ForwardCommandEvent(evt *base.CommandEvent) {
	json, err := json.Marshal(evt)
	if err != nil {
		fmt.Println("Error forwarding event to nats streaming server")
	}

	//fmt.Println(string(json))

	forwarder.stream.SnatNativeConnection.Publish(forwarder.stream.NatsPublishName, json)

	if forwarder.nextForwarder != nil {
		forwarder.nextForwarder.ForwardCommandEvent(evt)
	}
}
