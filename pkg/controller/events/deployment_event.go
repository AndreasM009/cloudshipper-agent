package events

// DeploymentCommand command
type DeploymentCommand struct {
	Name        string `json:"name"`
	Displayname string `json:"displayname"`
	Index       int    `json:"index"`
}

// DeploymentJob job
type DeploymentJob struct {
	Name        string              `json:"name"`
	Displayname string              `json:"displayname"`
	Commands    []DeploymentCommand `json:"commands"`
}

// DeploymentEvent event about definition
type DeploymentEvent struct {
	Event `json:"event"`
	Jobs  []DeploymentJob `json:"jobs"`
}
