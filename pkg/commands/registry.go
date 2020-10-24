package commands

import (
	"fmt"
)

// CommandRegistry the command registry
type CommandRegistry struct {
	Factories map[CommandType]*CommandFactories
}

// FromPropertiesHandler from props
type FromPropertiesHandler func(properties, inherited map[string]string) (interface{}, error)

// FromJSONHandler from yaml stream
type FromJSONHandler func(data []byte) (interface{}, error)

// CommandFactories a registered command's factory methods
type CommandFactories struct {
	FromProperties FromPropertiesHandler
	FromJSON       FromJSONHandler
}

var (
	theRegistry = CommandRegistry{
		Factories: make(map[CommandType]*CommandFactories),
	}
)

// GetRegistry gets registry
func GetRegistry() *CommandRegistry {
	return &theRegistry
}

// AddToCommandRegistry add a command to registry
func AddToCommandRegistry(register func(factories *CommandFactories) CommandType) bool {
	f := &CommandFactories{}
	t := register(f)
	theRegistry.Factories[t] = f
	return true
}

// CreateFromProps command from props
func (registry *CommandRegistry) CreateFromProps(commandType CommandType, props, inherited map[string]string) (interface{}, error) {
	if f, ok := registry.Factories[commandType]; ok {
		return f.FromProperties(props, inherited)
	}

	return nil, fmt.Errorf("Command type %d not found in registry", commandType)
}

// CreateFromJSON creates command from JSON stream
func (registry *CommandRegistry) CreateFromJSON(commandType CommandType, data []byte) (interface{}, error) {
	if f, ok := registry.Factories[commandType]; ok {
		return f.FromJSON(data)
	}

	return nil, fmt.Errorf("Command type %d not found in registry", commandType)
}
