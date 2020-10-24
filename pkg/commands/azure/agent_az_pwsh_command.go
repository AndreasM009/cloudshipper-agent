package azure

import (
	"encoding/json"
	"errors"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands"
)

const (
	pwshExecutable = "pwsh"
)

var (
	_ = commands.AddToCommandRegistry(register)
)

// AgentAzPwshCommand Powershell Core command
type AgentAzPwshCommand struct {
	AgentAzCommand `json:"azCommand"`
	ExecutableName string `json:"executableName"`
	ScriptToRun    string `json:"scriptToRun"`
	Arguments      string `json:"arguments"`
}

func register(factories *commands.CommandFactories) commands.CommandType {
	factories.FromProperties = newAzPwshCommandFromMap
	factories.FromJSON = newAzPwshCommandFromJSON
	return commands.AzPowershellCore
}

// NewAzPwshCommand create a new AzPwshCommand instance
func NewAzPwshCommand(scriptToRun, spname, spsecret, tenant, subscription, arguments, workingDirectory string) (*AgentAzPwshCommand, error) {
	if scriptToRun == "" {
		return nil, errors.New("scriptToRun can not be empty")
	}

	if scriptToRun == "" {
		return nil, errors.New("scriptToRun can not be empty")
	}

	if spname == "" {
		return nil, errors.New("spname can not be empty")
	}

	if spsecret == "" {
		return nil, errors.New("spsecret can not be empty")
	}

	if tenant == "" {
		return nil, errors.New("tenant can not be empty")
	}

	if subscription == "" {
		return nil, errors.New("subscription can not be empty")
	}

	return &AgentAzPwshCommand{
		AgentAzCommand: AgentAzCommand{
			Command: commands.Command{
				Type:             commands.AzPowershellCore,
				WorkingDirectory: workingDirectory,
			},
			ServicePrincipalName:   spname,
			ServicePrincipalSecret: spsecret,
			Tenant:                 tenant,
			Susbcription:           subscription,
		},
		ExecutableName: pwshExecutable,
		ScriptToRun:    scriptToRun,
		Arguments:      arguments,
	}, nil
}

// NewAzPwshCommandFromMap new instance from a map of properties
func newAzPwshCommandFromMap(props, inherited map[string]string) (interface{}, error) {
	var subscription, sp, secret, tenant, scriptToRun, arguments, workingdirectory string

	if val, ok := props["subscription"]; ok {
		subscription = val
	} else {
		return nil, errors.New("Azure Subscription not specified in definition")
	}

	if val, ok := props["serviceprincipal"]; ok {
		sp = val
	} else {
		return nil, errors.New("Azure Service Principal not specified in definition")
	}

	if val, ok := props["secret"]; ok {
		secret = val
	} else {
		return nil, errors.New("Azure ServicePrincipal secret not specified in definition")
	}

	if val, ok := props["tenant"]; ok {
		tenant = val
	} else {
		return nil, errors.New("Azure Tenant not specified in definition")
	}

	if val, ok := props["scriptToRun"]; ok {
		scriptToRun = val
	} else {
		return nil, errors.New("Azure Tenant not specified in definition")
	}

	if val, ok := props["arguments"]; ok {
		arguments = val
	}

	// inherited
	if val, ok := props["workingdirectory"]; ok {
		workingdirectory = val
	} else if val, ok := inherited["workingdirectory"]; ok {
		workingdirectory = val
	}

	pwsh, err := NewAzPwshCommand(scriptToRun, sp, secret, tenant, subscription, arguments, workingdirectory)
	if err != nil {
		return nil, err
	}

	return pwsh, nil
}

func newAzPwshCommandFromJSON(data []byte) (interface{}, error) {
	pwsh := &AgentAzPwshCommand{}

	if err := json.Unmarshal(data, pwsh); err != nil {
		return nil, err
	}

	return pwsh, nil
}
