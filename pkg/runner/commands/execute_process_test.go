package commands

import (
	"context"
	"testing"
	"time"

	fakeProxy "github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy/fake"
	"github.com/stretchr/testify/assert"
)

func TestExecutionUnknownExecutable(t *testing.T) {
	ctx := context.Background()
	proxy := fakeProxy.NewForReportingOnly(true)
	done := make(chan int, 1)

	result, err := executeProcessAsync(ctx, proxy, "unknown_executable", []string{}, []string{}, done)

	assert.NotNil(t, err)
	assert.Equal(t, 1, result)
	assert.Equal(t, 1, len(done))
}

func TestExecutionPwshExecutable(t *testing.T) {
	ctx := context.Background()
	proxy := fakeProxy.NewForReportingOnly(true)
	done := make(chan int, 1)

	result, err := executeProcessAsync(ctx, proxy, "pwsh", []string{"./execute-success.ps1"}, []string{}, done)

	assert.Nil(t, err)
	assert.Equal(t, 0, result)
	assert.Equal(t, 1, len(done))
}

func TestExecutionPwshExecutableError(t *testing.T) {
	ctx := context.Background()
	proxy := fakeProxy.NewForReportingOnly(true)
	done := make(chan int, 1)

	result, err := executeProcessAsync(ctx, proxy, "pwsh", []string{"./execute-error.ps1"}, []string{}, done)

	assert.NotNil(t, err)
	assert.Equal(t, 42, result)
	assert.Equal(t, 1, len(done))
}

func TestExecutionPwshExecutableCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	proxy := fakeProxy.NewForReportingOnly(true)
	done := make(chan int, 1)
	result := 0
	var err error = nil

	go func() {
		result, err = executeProcessAsync(ctx, proxy, "pwsh", []string{"./execute-cancel.ps1"}, []string{}, done)
	}()

	// wait and then cancel
	time.Sleep(time.Second * 3)
	cancel()

	timeoutctx, timeoutcancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer timeoutcancel()

	select {
	// timeout
	case <-timeoutctx.Done():
		assert.Fail(t, "test cancel process: timeout!")
	case result = <-done:
	}

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
	assert.Equal(t, 0, len(done))
}
