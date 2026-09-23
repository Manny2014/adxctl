package eventbus

// Original https://github.com/Manny2014/adxctl/blob/e70efca17f5fb25e09e1bb597e3371c1fdb4b9f4/internal/nats/runners/docker.go

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const (
	natsImage       = "nats:latest"
	containerName   = "adx-eventbus"
	natsPort        = "4222/tcp"
	startupFinished = "Server is ready"
)

func StartDockerRunner(ctx context.Context) error {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	reader, err := cli.ImagePull(ctx, natsImage, client.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull nats image: %w", err)
	}
	io.Copy(io.Discard, reader)
	reader.Close()

	port := network.MustParsePort(natsPort)

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: containerName,
		Config: &container.Config{
			Image: natsImage,
			Cmd: []string{
				"-js",
				"-m",
				"8222",
			},
			ExposedPorts: network.PortSet{
				port: struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			PortBindings: network.PortMap{
				port: []network.PortBinding{
					{
						// HostIP:   netip.IPv4Unspecified(),
						HostPort: "4222",
					},
					{
						// HostIP:   netip.IPv4Unspecified(),
						HostPort: "8222",
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if _, err := cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func StopDockerRunner(ctx context.Context) error {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	if _, err := cli.ContainerStop(ctx, containerName, client.ContainerStopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	if _, err := cli.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}
