package commands

import (
	"fmt"

	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"

	controllerCommand "github.com/andreasM009/cloudshipper-agent/pkg/commands"
)

var (
	theRegistry = CommandRegistry{
		Factories: make(map[controllerCommand.CommandType]*CommandFactories),
	}
)

// FactoryMethod to create concrete Command
type FactoryMethod func(controllerCommand interface{}, proxy proxy.ControllerProxy) (Command, error)

// CommandRegistry to register all Commands
type CommandRegistry struct {
	Factories map[controllerCommand.CommandType]*CommandFactories
}

// CommandFactories to register factory methods for command
type CommandFactories struct {
	Create FactoryMethod
}

// GetRegistry gets command registry
func GetRegistry() *CommandRegistry {
	return &theRegistry
}

// AddToCommandFactory adds command to registry
func AddToCommandFactory(f func(factories *CommandFactories) controllerCommand.CommandType) bool {
	factories := &CommandFactories{}
	cmdType := f(factories)
	theRegistry.Factories[cmdType] = factories
	return true
}

// CreateCommand creates concrete command for type 'commandType'
func (registry *CommandRegistry) CreateCommand(commandType controllerCommand.CommandType, agentCommand interface{}, proxy proxy.ControllerProxy) (Command, error) {
	if f, ok := registry.Factories[commandType]; ok {
		return f.Create(agentCommand, proxy)
	}

	return nil, fmt.Errorf("Unknown command for type: %s", commandType.String())
}
