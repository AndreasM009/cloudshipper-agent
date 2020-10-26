package commands

import "strings"

// CommandType specifies the type of command that is executed by the runner
type CommandType int

const (
	// UnknownCommand commond unknown
	UnknownCommand CommandType = iota
	// AzPowershellCore command
	AzPowershellCore
	// AzBash Command
	AzBash
	// AgentDownloadArtifacts command
	AgentDownloadArtifacts
)

func (c CommandType) String() string {
	return [...]string{
		"UnknownCommand",
		"AzPowershellCore",
		"AzBash",
		"AgentDownloadArtifacts",
	}[c]
}

// Parse CommandType from string
func Parse(s string) CommandType {
	s = strings.ToLower(s)
	switch s {
	case "unknowncommand":
		return UnknownCommand
	case "azpowershellcore":
		return AzPowershellCore
	case "azbash":
		return AzBash
	case "agentdownloadartifacts":
		return AgentDownloadArtifacts
	default:
		return UnknownCommand
	}
}
