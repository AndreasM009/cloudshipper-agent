package commands

import (
	"encoding/json"
	"errors"
)

var (
	_ = AddToCommandRegistry(registerDownload)
)

// AgentDownloadArtifactsCommand download artifacts
type AgentDownloadArtifactsCommand struct {
	Command      `json:"command"`
	ArtifactsURL string `json:"articafctsUrl"`
}

func registerDownload(factories *CommandFactories) CommandType {
	factories.FromProperties = newAgentnDownloadArtifactsCommandFromMap
	factories.FromJSON = newAgentnDownloadArtifactsCommandFromJSON
	return AgentDownloadArtifacts
}

// NewAgentDownloadArtifactsCommand new instance
func NewAgentDownloadArtifactsCommand(artifactsURL, workingdirectory string) (*AgentDownloadArtifactsCommand, error) {

	if artifactsURL == "" {
		return nil, errors.New("artifactsURL can not be empty")
	}

	return &AgentDownloadArtifactsCommand{
		Command: Command{
			Type:             AgentDownloadArtifacts,
			WorkingDirectory: workingdirectory,
		},
		ArtifactsURL: artifactsURL,
	}, nil
}

func newAgentnDownloadArtifactsCommandFromMap(props, inherited map[string]string) (interface{}, error) {
	var artifactsURL, workingDirectory string

	if val, ok := props["url"]; ok {
		artifactsURL = val
	}

	if val, ok := props["workingdirectory"]; ok {
		workingDirectory = val
	} else if val, ok := inherited["workingdirectory"]; ok {
		workingDirectory = val
	}

	return NewAgentDownloadArtifactsCommand(artifactsURL, workingDirectory)
}

func newAgentnDownloadArtifactsCommandFromJSON(data []byte) (interface{}, error) {
	cmd := &AgentDownloadArtifactsCommand{}
	if err := json.Unmarshal(data, cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}
