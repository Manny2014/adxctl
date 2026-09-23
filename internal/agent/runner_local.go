package agent

import (
	"adxctl/internal/config"
	"adxctl/pkg/adx"
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// LocalRunner executes tasks on the local machine.
type LocalRunner struct {
	cfg config.LocalRunnerConfig
}

// NewLocalRunner creates a new local runner.
func NewLocalRunner(cfg config.LocalRunnerConfig) (*LocalRunner, error) {
	return &LocalRunner{cfg: cfg}, nil
}

// Run executes the task locally.
func (r *LocalRunner) Run(ctx context.Context, subject string, task *adx.Task) error {
	fmt.Printf("LocalRunner: Running task %s of type %s on subject %s with runner type %s\n", task.ID, task.Type, subject, r.cfg.Type)

	switch r.cfg.Type {
	case "gemini-cli":
		fmt.Println("Executing with gemini-cli...")
		// Placeholder for gemini-cli logic
	case "antigravity-cli":
		fmt.Println("Executing with antigravity-cli...")
		// Placeholder for antigravity-cli logic
	case "local-script":
		if r.cfg.Script == "" {
			return fmt.Errorf("local-script runner requires a 'script' path in the configuration")
		}
		fmt.Printf("Executing local script: %s\n", r.cfg.Script)
		cmd := exec.CommandContext(ctx, r.cfg.Script)
		cmd.Stdin = bytes.NewReader(task.Payload) // Pass payload to script via stdin
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("local script execution failed: %w\nStderr: %s", err, stderr.String())
		}
		fmt.Printf("Local script output:\n%s\n", stdout.String())

	default:
		// a default behavior can be to just print the task
		fmt.Printf("Received task: %+v\n", task)
	}

	return nil
}
