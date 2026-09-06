package agent

import "context"

// Runner defines the interface for executing a task in a specific environment.
type Runner interface {
	Run(ctx context.Context, task *Task) error
}

// Task represents a unit of work to be processed by an agent.
type Task struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Payload     []byte      `json:"payload"`
	Complexity  Complexity  `json:"complexity"`
}

// Complexity represents the assessed complexity of a task.
type Complexity string

const (
	ComplexitySimple   Complexity = "simple"
	ComplexityComplex  Complexity = "complex"
	ComplexityUnknown  Complexity = "unknown"
)
