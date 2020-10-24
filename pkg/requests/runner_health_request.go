package requests

// RunnerHealthRequest to check if agent is healthy or not
type RunnerHealthRequest struct {
	Request   `json:"request"`
	IsHealthy bool `json:"ishealthy"`
}

// NewRunnerHealthRequest new instance
func NewRunnerHealthRequest() *RunnerHealthRequest {
	return &RunnerHealthRequest{
		Request: Request{
			RequestType: RunnerHealth,
		},
		IsHealthy: false,
	}
}
