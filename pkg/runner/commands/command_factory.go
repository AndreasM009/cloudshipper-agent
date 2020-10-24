package commands

import (
	"encoding/json"

	controllerCommand "github.com/andreasM009/cloudshipper-agent/pkg/commands"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
)

// CreateCommand create command from stream
func CreateCommand(data []byte, proxy proxy.ControllerProxy) (Command, error) {
	carrier := controllerCommand.CommandCarrier{}

	if err := json.Unmarshal(data, &carrier); err != nil {
		return nil, err
	}

	stream, err := json.Marshal(carrier.Data)
	if err != nil {
		return nil, err
	}

	cntrlcmd, err := controllerCommand.GetRegistry().CreateFromJSON(carrier.CarrierForType, stream)
	if err != nil {
		return nil, err
	}

	cmd, err := GetRegistry().CreateCommand(carrier.CarrierForType, cntrlcmd, proxy)
	if err != nil {
		return nil, err
	}

	return cmd, err
}
