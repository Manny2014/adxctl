package eventbus

import (
	"os/exec"
	"testing"
)

func TestLocalRunner(t *testing.T) {
	// This test requires nats-server to be in the PATH.
	if _, err := exec.LookPath("nats-server"); err != nil {
		t.Skip("Skipping local runner tests: nats-server not found in PATH")
	}

	pid, err := StartLocalRunner()
	if err != nil {
		t.Fatalf("StartLocalRunner() error = %v", err)
	}

	if pid == 0 {
		t.Error("StartLocalRunner() returned PID 0")
	}

	err = StopLocalRunner(pid)
	if err != nil {
		t.Fatalf("StopLocalRunner() error = %v", err)
	}
}
