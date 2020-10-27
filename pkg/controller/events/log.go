package events

import "github.com/andreasM009/cloudshipper-agent/pkg/logs"

// Log entry
type Log struct {
	IsWarning bool   `json:"isWarning"`
	IsError   bool   `json:"isError"`
	IsInfo    bool   `json:"isInfo"`
	Message   string `json:"message"`
}

// NewLogFromLogs log
func NewLogFromLogs(l *logs.LogMessage) Log {
	res := Log{
		IsWarning: false,
		IsError:   false,
		IsInfo:    false,
		Message:   l.Message,
	}

	if l.LogType == logs.Error {
		res.IsError = true
		return res
	} else if l.LogType == logs.Warning {
		res.IsWarning = true
		return res
	} else {
		res.IsInfo = true
		return res
	}
}
