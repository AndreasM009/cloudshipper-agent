package proxy

import (
	"github.com/andreasM009/cloudshipper-agent/pkg/logs"
)

// ControllerProxy communicates with controller
type ControllerProxy interface {
	Report(logs.LogMessage) error
	ReportError(exitcode int) error
	GetAgentCommand() ([]byte, error)
}
