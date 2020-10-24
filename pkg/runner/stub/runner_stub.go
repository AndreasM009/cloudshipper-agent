package stub

import "context"

// HealthyHandler to handle health
type HealthyHandler func() bool

// CancellationHandler to handle cancellation
type CancellationHandler func() bool

// RunnerStub stub interface
type RunnerStub interface {
	SetHandler(healthy HealthyHandler, cancellation CancellationHandler) error
	StartListeningAndBlock(ctx context.Context, commandrunner <-chan int) error
}
