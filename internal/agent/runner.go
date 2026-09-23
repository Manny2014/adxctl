package agent

import (
	"adxctl/pkg/adx"
	"context"
)

// Runner defines the interface for executing a task in a specific environment.
type Runner interface {
	Run(ctx context.Context, subject string, task *adx.Task) error
}
