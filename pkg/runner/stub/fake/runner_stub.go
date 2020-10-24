package fake

import (
	"context"
	"fmt"

	"github.com/andreasM009/cloudshipper-agent/pkg/runner/stub"
)

// RunnerStub fake runner stub
type RunnerStub struct {
	healthyHandler      stub.HealthyHandler
	cancellationHandler stub.CancellationHandler
}

// NewFakeRunnerStub new instance
func NewFakeRunnerStub(healthy stub.HealthyHandler, cancellation stub.CancellationHandler) *RunnerStub {
	return &RunnerStub{
		healthyHandler:      healthy,
		cancellationHandler: cancellation,
	}
}

// Cancel cancels runner
func (runner *RunnerStub) Cancel() {
	runner.cancellationHandler()
}

// IsHealthy is runner healthy
func (runner *RunnerStub) IsHealthy() bool {
	return runner.healthyHandler()
}

// SetHandler sets callback handler
func (runner *RunnerStub) SetHandler(healthy stub.HealthyHandler, cancellation stub.CancellationHandler) error {
	return nil
}

// StartListeningAndBlock starts listening and blocks until command runner stops
func (runner *RunnerStub) StartListeningAndBlock(ctx context.Context, commandrunner <-chan int) error {

	exitcode := <-commandrunner
	fmt.Println(fmt.Sprintf("Runner stpped with exitcode: %d", exitcode))
	return nil
}
