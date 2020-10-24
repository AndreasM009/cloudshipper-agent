package fake

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/logs"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
)

const (
	executableName         = "EXECUTABLE_NAME"
	scriptToRun            = "SCRIPT_TO_RUN"
	servicePrincipalName   = "SERVICEPRINCIPAL_NAME"
	servicePrincipalSecret = "SERVICEPRINCIPAL_SECRET"
	tenant                 = "TENANT"
	subscription           = "SUBSCRIPTION"
)

// ControllerProxy fake the proxy
type ControllerProxy struct {
	executableName         string
	scriptToRun            string
	servicePrincipalName   string
	servicePrincipalSecret string
	tenant                 string
	subscription           string
	artifactsURL           string
	numOfCommandToRun      int
	numOfCommandsRan       int
	silentReporting        bool
}

// NewForReportingOnly simple proxy that logs to stdout and stderr
func NewForReportingOnly(silent bool) proxy.ControllerProxy {
	return &ControllerProxy{
		numOfCommandToRun: 0,
		numOfCommandsRan:  0,
		silentReporting:   silent,
	}
}

// NewFromEnvironment creates a new instance from ENV variables
func NewFromEnvironment(numOfCommandsToRun int) proxy.ControllerProxy {
	return &ControllerProxy{
		executableName:         os.Getenv(executableName),
		scriptToRun:            os.Getenv(scriptToRun),
		servicePrincipalName:   os.Getenv(servicePrincipalName),
		servicePrincipalSecret: os.Getenv(servicePrincipalSecret),
		tenant:                 os.Getenv(tenant),
		subscription:           os.Getenv(subscription),
		artifactsURL:           "",
		numOfCommandToRun:      numOfCommandsToRun,
		numOfCommandsRan:       0,
		silentReporting:        false,
	}
}

// Report implements ControllerProxy
func (p *ControllerProxy) Report(l logs.LogMessage) error {
	if p.silentReporting {
		return nil
	}

	if l.LogType == logs.Error {
		fmt.Println(fmt.Sprintf("[ERROR]: %s", l.Message))
	} else if l.LogType == logs.Info {
		fmt.Println(fmt.Sprintf("[INFO]: %s", l.Message))
	} else {
		fmt.Println(fmt.Sprintf("[WARNING]: %s", l.Message))
	}

	return nil
}

// ReportError error report during command execution
func (p *ControllerProxy) ReportError(exitcide int) error {
	fmt.Println(fmt.Sprintf("Error in command: exitcode %d", exitcide))
	return nil
}

// GetAgentCommand implements ControllerProxy
func (p *ControllerProxy) GetAgentCommand() ([]byte, error) {
	if p.numOfCommandsRan == p.numOfCommandToRun {
		return nil, nil
	}

	cmd, err := azure.NewAzPwshCommand(p.scriptToRun,
		p.servicePrincipalName, p.servicePrincipalSecret, p.tenant, p.subscription, "", "")

	if err != nil {
		log.Panic(err)
	}

	carrier := commands.CommandCarrier{
		CarrierForType: cmd.Type,
		Data:           &cmd,
	}

	buffer, err := json.Marshal(carrier)

	if err != nil {
		log.Panic(err)
	}

	str := string(buffer)
	fmt.Println(str)

	p.numOfCommandsRan = 1
	return buffer, nil
}
