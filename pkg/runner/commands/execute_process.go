package commands

import (
	"context"
	"os/exec"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
)

func executeProcessAsync(
	ctx context.Context, proxy proxy.ControllerProxy,
	exectutable string, args []string, env []string, done chan int) (int, error) {
	// look for executable
	ps, err := exec.LookPath(exectutable)
	if err != nil {
		proxy.Report(logs.NewErrorLog(err.Error()))
		done <- 1
		return 1, err
	}

	// create command
	pcmd := exec.Command(ps, args...)

	// set Stdout
	pcmd.Stdout = &stdOut2LogsWriter{
		proxy: proxy,
	}

	// set Stderr
	pcmd.Stderr = &stdErr2LogsWriter{
		proxy: proxy,
	}

	// execute process
	err = pcmd.Start()
	if err != nil {
		proxy.Report(logs.NewErrorLog(err.Error()))
		done <- 1
		return 1, err
	}

	// create result channel
	var resulterr error = nil
	resultchan := make(chan int)

	// execute command async
	go func() {
		if err := pcmd.Wait(); err != nil {
			proxy.Report(logs.NewErrorLog(err.Error()))
			resulterr = err
			resultchan <- pcmd.ProcessState.ExitCode()
		} else {
			resultchan <- pcmd.ProcessState.ExitCode()
		}
	}()

	exitcode := 0

	// wait for process result or for cancellation
	select {
	// canceled
	case <-ctx.Done():
		pcmd.Process.Kill()     // kill the running process
		exitcode = <-resultchan // get exitcode of process
		resulterr = nil
	// proccess finished
	case exitcode = <-resultchan:
	}

	// notify done channel
	done <- exitcode
	// return exitcode and error
	return exitcode, resulterr
}
