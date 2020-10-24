package events

// CommandEvent event
type CommandEvent struct {
	Event
	JobName            string
	JobDisplayName     string
	CommandName        string
	CommandDisplayName string
	CommandIndex       int
	Logs               []Log
}