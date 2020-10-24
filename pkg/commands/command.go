package commands

// Command base command
type Command struct {
	Type             CommandType `json:"type"`
	WorkingDirectory string      `json:"workingDirectory"`
}
