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
	Event
	Jobs     []DeploymentJob `json:"jobs"`
	Started  bool            `json:"started"`
	Finished bool            `json:"finished"`
	Exitcode int             `json:"exitcode"`
}
