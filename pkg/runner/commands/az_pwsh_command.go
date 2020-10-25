package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"

	agentcommands "github.com/andreasM009/cloudshipper-agent/pkg/commands"
	"github.com/andreasM009/cloudshipper-agent/pkg/commands/azure"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/settings"
)

var (
	_ = AddToCommandFactory(registerAzPwshCommand)
)

// AzPwshCommand executes a PowerShell Core Script
type AzPwshCommand struct {
	controllerCommand *azure.AgentAzPwshCommand
	controllerProxy   proxy.ControllerProxy
	done              chan int
}

func registerAzPwshCommand(factories *CommandFactories) agentcommands.CommandType {
	factories.Create = newAzPwshCommand
	return agentcommands.AzPowershellCore
}

// newAzPwshCommand creates an Azure PowerShell Core command
func newAzPwshCommand(cmd interface{}, proxy proxy.ControllerProxy) (Command, error) {
	if c, ok := cmd.(*azure.AgentAzPwshCommand); ok {
		return &AzPwshCommand{
			controllerCommand: c,
			controllerProxy:   proxy,
			done:              make(chan int, 1),
		}, nil
	}
	return nil, errors.New("type assertion failed")
}

// Execute runs command
func (cmd *AzPwshCommand) Execute(ctx context.Context) (int, error) {
	artifactsDirectory := settings.GetArtifactsDirectory()

	args := []string{
		"./azpwsh.ps1",
		"-ScriptToRun",
		cmd.controllerCommand.ScriptToRun,
		"-ArgumentsToRun",
		cmd.controllerCommand.Arguments,
		"-Sp",
		cmd.controllerCommand.ServicePrincipalName,
		"-Secret",
		cmd.controllerCommand.ServicePrincipalSecret,
		"-Tenant",
		cmd.controllerCommand.Tenant,
		"-Subscription",
		cmd.controllerCommand.Susbcription,
		"-ArtifactsDirectory",
		artifactsDirectory,
	}

	if cmd.controllerCommand.WorkingDirectory != "" {
		args = append(args, "-WorkingDirectory", cmd.controllerCommand.WorkingDirectory)
	}

	env := []string{
		fmt.Sprintf("AD_SERVICEPRINCIPAL_NAME=%s", cmd.controllerCommand.ServicePrincipalName),
		fmt.Sprintf("AD_SERVICEPRINCIPAL_SECRET=%s", cmd.controllerCommand.ServicePrincipalSecret),
		fmt.Sprintf("AD_TENANT=%s", cmd.controllerCommand.Tenant),
		fmt.Sprintf("AZURE_SUBSCRIPTION=%s", cmd.controllerCommand.Susbcription),
		"PSModulePath=/usr/local/share/powershell/Modules:/usr/local/microsoft/powershell/6/Modules",
		fmt.Sprintf("HOME=/Users/%s", os.Getenv("USER")),
	}

	if err := cmd.controllerProxy.Report(logs.NewInfoLog("Executing Azure PowerShell.")); err != nil {
		fmt.Println(err)
	}

	return executeProcessAsync(ctx, cmd.controllerProxy, cmd.controllerCommand.ExecutableName, args, env, cmd.done)
}

// ExecuteAsync executes command asynchronously
func (cmd *AzPwshCommand) ExecuteAsync(ctx context.Context) error {
	// todo: validate command before starting go routine
	go cmd.Execute(ctx)

	return nil
}

// Done channel
func (cmd *AzPwshCommand) Done() <-chan int {
	return cmd.done
}
