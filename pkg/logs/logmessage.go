package logs

// LogType type of log
type LogType int

const (
	// Info log contains information
	Info LogType = 0
	// Warning log contains warning
	Warning LogType = 1
	// Error log contains error
	Error LogType = 2
)

// LogMessage to report logs to controller
type LogMessage struct {
	LogType LogType
	Message string
}

// NewErrorLog create a new log of type error
func NewErrorLog(msg string) LogMessage {
	return LogMessage{
		LogType: Error,
		Message: msg,
	}
}

// NewInfoLog create a new log of type info
func NewInfoLog(msg string) LogMessage {
	return LogMessage{
		LogType: Info,
		Message: msg,
	}
}

// NewWarningLog create a new log of type warning
func NewWarningLog(msg string) LogMessage {
	return LogMessage{
		LogType: Warning,
		Message: msg,
	}
}
