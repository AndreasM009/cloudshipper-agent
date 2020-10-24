package requests

// ControllerQueryRunnerRequest Request to get concrete Request for agent runner
type ControllerQueryRunnerRequest struct {
	Request `json:"request"`
}

// NewControllerQueryRunnerRequest new instance
func NewControllerQueryRunnerRequest() *ControllerQueryRunnerRequest {
	return &ControllerQueryRunnerRequest{
		Request: Request{
			RequestType: ControllerQueryCommand,
		},
	}
}
