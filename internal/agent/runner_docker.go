package agent

import (
	"context"
	"fmt"
)

// DockerRunner executes tasks in a Docker container.
type DockerRunner struct{}

// NewDockerRunner creates a new Docker runner.
func NewDockerRunner() (*DockerRunner, error) {
	return &DockerRunner{}, nil
}

// Run executes the task in a Docker container.
func (r *DockerRunner) Run(ctx context.Context, task *Task) error {
	fmt.Printf("DockerRunner: Running task %s of type %s\n", task.ID, task.Type)
	// Placeholder for Docker task execution logic
	return nil
}
