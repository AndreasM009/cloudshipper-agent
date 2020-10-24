package requests

// RequestType types
type RequestType int

const (
	// UnknownRequest unknown
	UnknownRequest RequestType = iota
	// RunnerCancel cancel runner
	RunnerCancel
	// RunnerHealth is healthy
	RunnerHealth
	// ControllerLog log
	ControllerLog
	// ControllerQueryCommand wuery next command
	ControllerQueryCommand
	// ControllerReportCommandError error in command execution
	ControllerReportCommandError
)
