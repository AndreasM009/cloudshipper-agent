package events

// CommandEvent event
type CommandEvent struct {
	Event
	JobName            string `json:"jobName"`
	JobDisplayName     string `json:"jobDisplayName"`
	CommandName        string `json:"commandName"`
	CommandDisplayName string `json:"commandDisplayName"`
	CommandIndex       int    `json:"commandIndex"`
	Logs               []Log  `json:"logs"`
}
