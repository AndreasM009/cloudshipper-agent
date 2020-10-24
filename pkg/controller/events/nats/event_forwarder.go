package nats

import (
	"encoding/json"
	"fmt"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	base "github.com/andreasM009/cloudshipper-agent/pkg/controller/events"
)

// EventForwarder forwards events to nats stream
type EventForwarder struct {
	stream *channel.NatsStreamingChannel
}

// NewNatsStreamEventForwarder new instance
func NewNatsStreamEventForwarder(stream *channel.NatsStreamingChannel) base.EventForwarder {
	return &EventForwarder{
		stream: stream,
	}
}

// ForwardDeploymentEvent event
func (forwarder *EventForwarder) ForwardDeploymentEvent(evt *base.DeploymentEvent) {
	json, err := json.Marshal(evt)
	if err != nil {
		fmt.Println("Error forwarding event to nats streaming server")
	}

	forwarder.stream.SnatConnection.Publish(forwarder.stream.NatsPublishName, json)
}

// ForwardCommandEvent event
func (forwarder *EventForwarder) ForwardCommandEvent(evt *base.CommandEvent) {
	json, err := json.Marshal(evt)
	if err != nil {
		fmt.Println("Error forwarding event to nats streaming server")
	}

	forwarder.stream.SnatConnection.Publish(forwarder.stream.NatsPublishName, json)
}
