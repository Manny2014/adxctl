package eventbus

import (
	"fmt"
	"os"
	"os/exec"

	"syscall"
)

func StartLocalRunner() (int, error) {
	cmd := exec.Command("nats-server", "-js", "-p", "4222")
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start nats-server: %w", err)
	}

	return cmd.Process.Pid, nil
}

func StopLocalRunner(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with pid %d: %w", pid, err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to terminate process with pid %d: %w", pid, err)
	}
	return nil
}
