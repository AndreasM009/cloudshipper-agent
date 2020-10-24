package definition

import (
	"gopkg.in/yaml.v2"
)

// Step step
type Step struct {
	Command          string
	Displayname      string
	WorkingDirectory string `yaml:"working-directory"`
	With             map[string]string
}

// Job definition
type Job struct {
	Displayname      string
	Steps            []Step
	WorkingDirectory string `yaml:"working-directory"`
}

// Deployment definition
type Deployment struct {
	Jobs map[string]Job
}

// NewFromYaml new instance from yaml
func NewFromYaml(data []byte) (*Deployment, error) {
	deployment := &Deployment{}

	err := yaml.Unmarshal(data, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}
