package nats

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	"github.com/andreasM009/cloudshipper-agent/pkg/requests"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
)

// ControllerProxy implements ControllerProxy based on nats messages
type ControllerProxy struct {
	natsChannel *channel.NatsChannel
}

// NewNatsControllerProxy new instance
func NewNatsControllerProxy(channel *channel.NatsChannel) (proxy.ControllerProxy, error) {
	if nil == channel {
		return nil, errors.New("Channel can not be nil")
	}
	proxy := &ControllerProxy{
		natsChannel: channel,
	}
	return proxy, nil
}

// Report logs
func (proxy *ControllerProxy) Report(log logs.LogMessage) error {
	carrier := requests.RequestCarrier{
		CarrierForType: requests.ControllerLog,
		Data:           requests.NewControllerLogRequest(log),
	}

	json, err := json.Marshal(carrier)
	if err != nil {
		return err
	}

	if err := proxy.natsChannel.NatsNativeConn.Publish(proxy.natsChannel.NatsPublishName, json); err != nil {
		return err
	}

	return nil
}

// ReportError error in command execution
func (proxy *ControllerProxy) ReportError(exitcode int) error {
	carrier := requests.RequestCarrier{
		CarrierForType: requests.ControllerReportCommandError,
		Data:           requests.NewControllerReportErrorRequest(exitcode),
	}

	json, err := json.Marshal(carrier)
	if err != nil {
		return err
	}

	if err := proxy.natsChannel.NatsNativeConn.Publish(proxy.natsChannel.NatsPublishName, json); err != nil {
		return err
	}

	return nil
}

// GetAgentCommand gets new command to execute
func (proxy *ControllerProxy) GetAgentCommand() ([]byte, error) {
	carrier := requests.RequestCarrier{
		CarrierForType: requests.ControllerQueryCommand,
		Data:           requests.NewControllerQueryRunnerRequest(),
	}

	json, err := json.Marshal(carrier)
	if err != nil {
		return nil, err
	}

	msg, err := proxy.natsChannel.NatsNativeConn.Request(proxy.natsChannel.NatsPublishName, json, 60*time.Second)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
