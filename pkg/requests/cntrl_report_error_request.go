package requests

// ControllerReportErrorRequest error report
type ControllerReportErrorRequest struct {
	Request  `json:"request"`
	Exitcode int `json:"exitcode"`
}

// NewControllerReportErrorRequest new instance
func NewControllerReportErrorRequest(exitcode int) *ControllerReportErrorRequest {
	return &ControllerReportErrorRequest{
		Request: Request{
			RequestType: ControllerReportCommandError,
		},
		Exitcode: exitcode,
	}
}
