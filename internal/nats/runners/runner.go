package nats

import (
	"context"
	"fmt"
	"strings"

	"adxctl/internal/config"
)

// StatusInfo represents the runtime status of a NATS server instance
type StatusInfo struct {
	Running       bool   `json:"running"`
	State         string `json:"state"`
	Runner        string `json:"runner"`
	ConnectURL    string `json:"connect_url"`
	MonitorURL    string `json:"monitor_url"`
	Port          int    `json:"port"`
	PID           int    `json:"pid,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
	ContainerID   string `json:"container_id,omitempty"`
	Image         string   `json:"image,omitempty"`
	StoreDir      string   `json:"store_dir,omitempty"`
	Logs          []string `json:"logs,omitempty"`
}

// Runner defines the interface for running a NATS server instance
type Runner interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (*StatusInfo, error)
}

// Options contains parameters passed to configure the NATS runner
type Options struct {
	Port        int
	StoreDir    string
	DockerImage string
}

// New creates a new Runner based on the runner type string
func New(runnerType string, opts Options) (Runner, error) {
	if opts.Port <= 0 {
		opts.Port = config.DefaultPort
	}
	if opts.DockerImage == "" {
		opts.DockerImage = config.DefaultDockerImage
	}

	switch strings.ToLower(runnerType) {
	case "local-cli", "local":
		return NewLocalRunner(opts), nil
	case "docker":
		return NewDockerRunner(opts), nil
	default:
		return nil, fmt.Errorf("unsupported runner %q (supported: %q, %q)", runnerType, config.RunnerLocalCLI, config.RunnerDocker)
	}
}
