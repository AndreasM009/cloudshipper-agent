package commands

import (
	"context"
	"errors"

	controllerCommand "github.com/andreasM009/cloudshipper-agent/pkg/commands"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/settings"
)

var (
	_ = AddToCommandFactory(registerDownloadArtifactsCommand)
)

// DownloadArtifactsCommand command
type DownloadArtifactsCommand struct {
	controllerCommand *controllerCommand.AgentDownloadArtifactsCommand
	controllerProxy   proxy.ControllerProxy
	done              chan int
}

func registerDownloadArtifactsCommand(factories *CommandFactories) controllerCommand.CommandType {
	factories.Create = newDownloadArtifactsCommand
	return controllerCommand.AgentDownloadArtifacts
}

// newDownloadArtifactsCommand creates an Azure PowerShell Core command
func newDownloadArtifactsCommand(cmd interface{}, proxy proxy.ControllerProxy) (Command, error) {
	if c, ok := cmd.(*controllerCommand.AgentDownloadArtifactsCommand); ok {
		return &DownloadArtifactsCommand{
			controllerCommand: c,
			controllerProxy:   proxy,
			done:              make(chan int, 1),
		}, nil
	}

	return nil, errors.New("type assertion failed")
}

// Execute execute
func (cmd *DownloadArtifactsCommand) Execute(ctx context.Context) (int, error) {
	artifactsDirectory := settings.GetArtifactsDirectory()

	args := []string{
		"-P",
		artifactsDirectory,
		cmd.controllerCommand.ArtifactsURL,
	}

	env := []string{}

	return executeProcessAsync(ctx, cmd.controllerProxy, "wget", args, env, cmd.done)
}

// ExecuteAsync execute
func (cmd *DownloadArtifactsCommand) ExecuteAsync(ctx context.Context) error {
	// todo: validate command before starting go routine
	go cmd.Execute(ctx)

	return nil
}

// Done done channel
func (cmd *DownloadArtifactsCommand) Done() <-chan int {
	return cmd.done
}
