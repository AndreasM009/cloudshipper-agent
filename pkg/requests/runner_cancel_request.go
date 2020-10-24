package requests

// RunnerCancelRequest to send cancellation request to runner
type RunnerCancelRequest struct {
	Request `json:"request"`
}

// NewRunnerCancelRequest new instance
func NewRunnerCancelRequest() *RunnerCancelRequest {
	return &RunnerCancelRequest{
		Request: Request{
			RequestType: RunnerCancel,
		},
	}
}
