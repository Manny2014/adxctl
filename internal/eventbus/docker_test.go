package eventbus

import (
	"context"
	"testing"
)

func TestDockerRunner(t *testing.T) {
	// These tests require a running docker daemon and are more of an integration test.
	// In a real scenario, you would mock the docker client.
	t.Skip("Skipping docker runner tests in CI")

	ctx := context.Background()
	err := StartDockerRunner(ctx)
	if err != nil {
		t.Fatalf("StartDockerRunner() error = %v", err)
	}

	// Here you would add assertions to check if the container is running.

	err = StopDockerRunner(ctx)
	if err != nil {
		t.Fatalf("StopDockerRunner() error = %v", err)
	}

	// Here you would add assertions to check if the container is stopped and removed.
}
