package agent

import (
	"context"
)

// Orchestrator is responsible for consuming events and dispatching tasks to the appropriate runners.
type Orchestrator struct {
	// dependencies like eventbus client, AI client, and runners will be added here
}

// NewOrchestrator creates a new agent orchestrator.
func NewOrchestrator() (*Orchestrator, error) {
	return &Orchestrator{}, nil
}

// Run starts the orchestrator's event consumption loop.
func (o *Orchestrator) Run(ctx context.Context) error {
	// Subscription and task dispatching logic will go here.
	<-ctx.Done()
	return nil
}
