package agent

import (
	"adxctl/pkg/adx"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// DockerRunner executes tasks in a Docker container.
type DockerRunner struct {
	cli       *client.Client
	ImageName string
}

// NewDockerRunner creates a new Docker runner.
func NewDockerRunner(imageName string) (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	if imageName == "" {
		imageName = "adxctl-worker"
	}
	return &DockerRunner{cli: cli, ImageName: imageName}, nil
}

// Run executes the task in a Docker container.
func (r *DockerRunner) Run(ctx context.Context, subject string, task *adx.Task) error {
	fmt.Printf("DockerRunner: Running task %s of type %s on subject %s\n", task.ID, task.Type, subject)

	// Build the image
	buildOptions := client.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{r.ImageName},
		Remove:     true,
	}
	buildContext, err := os.Open(".")
	if err != nil {
		return fmt.Errorf("failed to open build context: %w", err)
	}
	defer buildContext.Close()

	buildResp, err := r.cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer buildResp.Body.Close()

	// Wait for the image build to complete
	_, err = io.Copy(os.Stdout, buildResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read build response: %w", err)
	}

	// Create the container
	resp, err := r.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: r.ImageName,
			Cmd:   []string{"worker", "run", "--subject", subject, "--runtime", "local"},
			Env: []string{
				fmt.Sprintf("NATS_URL=%s", os.Getenv("NATS_URL")),
			},
		},
		HostConfig:       &container.HostConfig{},
		NetworkingConfig: &network.NetworkingConfig{},
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if _, err := r.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("DockerRunner: Started container %s for task %s\n", resp.ID, task.ID)

	return nil
}
