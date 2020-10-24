package requests

import (
	"github.com/andreasM009/cloudshipper-agent/pkg/logs"
)

// ControllerLogRequest Request
type ControllerLogRequest struct {
	Request `json:"request"`
	Log     logs.LogMessage `json:"log"`
}

// NewControllerLogRequest new instance
func NewControllerLogRequest(log logs.LogMessage) *ControllerLogRequest {
	return &ControllerLogRequest{
		Request: Request{
			RequestType: ControllerLog,
		},
		Log: log,
	}
}
