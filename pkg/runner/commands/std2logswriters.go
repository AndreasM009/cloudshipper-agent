package commands

import (
	"github.com/andreasM009/cloudshipper-agent/pkg/logs"

	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
)

type stdErr2LogsWriter struct {
	proxy proxy.ControllerProxy
}

type stdOut2LogsWriter struct {
	proxy proxy.ControllerProxy
}

func (w *stdErr2LogsWriter) Write(p []byte) (n int, err error) {
	w.proxy.Report(logs.NewErrorLog(string(p)))
	return len(p), nil
}

func (w *stdOut2LogsWriter) Write(p []byte) (n int, err error) {
	w.proxy.Report(logs.NewInfoLog(string(p)))
	return len(p), nil
}
