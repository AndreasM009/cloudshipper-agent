package commands

import (
	"context"
)

// Command interface
type Command interface {
	// Execute command
	// @return 0 == success
	// @return -1 == command was canceled
	// @return >0 == error code
	Execute(ctx context.Context) (int, error)

	// ExecuteAsync rund commmand asynchronously
	ExecuteAsync(ctx context.Context) error

	// Done channel
	// @return 0 == success
	// @return -1 == command was canceled
	// @return >0 == error code
	Done() <-chan int
}
