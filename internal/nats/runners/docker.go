package nats

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultContainerName = "adxctl-nats"

// DockerRunner manages a nats-server container running under Docker
type DockerRunner struct {
	opts          Options
	containerName string
}

// NewDockerRunner constructs a DockerRunner
func NewDockerRunner(opts Options) *DockerRunner {
	return &DockerRunner{
		opts:          opts,
		containerName: defaultContainerName,
	}
}

func (r *DockerRunner) ensureDockerCLI() (string, error) {
	bin, err := exec.LookPath("docker")
	if err != nil {
		return "", fmt.Errorf("docker executable not found in PATH. Please install Docker or start the Docker daemon")
	}
	return bin, nil
}

// Status checks the status of the Docker container
func (r *DockerRunner) Status(ctx context.Context) (*StatusInfo, error) {
	info := &StatusInfo{
		Runner:        "docker",
		Port:          r.opts.Port,
		ConnectURL:    fmt.Sprintf("nats://localhost:%d", r.opts.Port),
		MonitorURL:    "http://localhost:8222",
		ContainerName: r.containerName,
		Image:         r.opts.DockerImage,
		StoreDir:      r.opts.StoreDir,
		State:         "not created",
		Running:       false,
	}

	dockerBin, err := r.ensureDockerCLI()
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, dockerBin, "inspect", "--format", "{{.State.Status}}|{{.Id}}|{{.Config.Image}}", r.containerName)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// inspect returns non-zero if container does not exist
		return info, nil
	}

	parts := strings.Split(strings.TrimSpace(out.String()), "|")
	if len(parts) > 0 && parts[0] != "" {
		info.State = parts[0]
		info.Running = (parts[0] == "running")
	}
	if len(parts) > 1 && parts[1] != "" {
		cid := parts[1]
		if len(cid) > 12 {
			cid = cid[:12]
		}
		info.ContainerID = cid
	}
	if len(parts) > 2 && parts[2] != "" {
		info.Image = parts[2]
	}

	// Attach last 5 log lines when container is running
	if info.Running {
		logsCmd := exec.CommandContext(ctx, dockerBin, "logs", "--tail", "5", r.containerName)
		var logsOut bytes.Buffer
		logsCmd.Stdout = &logsOut
		logsCmd.Stderr = &logsOut // NATS writes to stderr
		if logsCmd.Run() == nil {
			for _, line := range strings.Split(strings.TrimSpace(logsOut.String()), "\n") {
				if line != "" {
					info.Logs = append(info.Logs, line)
				}
			}
		}
	}

	return info, nil
}

// Start runs or starts the nats-server container
func (r *DockerRunner) Start(ctx context.Context) error {
	dockerBin, err := r.ensureDockerCLI()
	if err != nil {
		return err
	}

	st, err := r.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}

	if st.Running {
		return fmt.Errorf("docker container %q is already running", r.containerName)
	}

	if st.State == "exited" || st.State == "created" {
		// Restart existing stopped container
		startCmd := exec.CommandContext(ctx, dockerBin, "start", r.containerName)
		var stderr bytes.Buffer
		startCmd.Stderr = &stderr
		if err := startCmd.Run(); err != nil {
			return fmt.Errorf("failed to start existing container %q: %s", r.containerName, stderr.String())
		}
		fmt.Printf("Started existing container %q with image %q\n", r.containerName, r.opts.DockerImage)
		fmt.Printf("Connect URL: nats://localhost:%d\n", r.opts.Port)
		fmt.Printf("Monitoring:  http://localhost:8222\n")
		return nil
	}

	portMapping := fmt.Sprintf("%s:4222", strconv.Itoa(r.opts.Port))
	args := []string{
		"run",
		"-d",
		"--name", r.containerName,
		"-p", portMapping,
		"-p", "8222:8222",
	}

	var absStoreDir string
	if r.opts.StoreDir != "" {
		var err error
		absStoreDir, err = filepath.Abs(r.opts.StoreDir)
		if err != nil {
			return fmt.Errorf("invalid store_dir path %q: %w", r.opts.StoreDir, err)
		}
		if err := os.MkdirAll(absStoreDir, 0755); err != nil {
			return fmt.Errorf("failed to create store_dir directory %q: %w", absStoreDir, err)
		}
		args = append(args, "-v", fmt.Sprintf("%s:/data", absStoreDir))
	}

	args = append(args, r.opts.DockerImage)

	// Append jetstream
	args = append(args, "-js")

	if absStoreDir != "" {
		args = append(args, "-sd", "/data")
	}

	runCmd := exec.CommandContext(ctx, dockerBin, args...)
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("docker run failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}

	containerID := strings.TrimSpace(stdout.String())
	if len(containerID) > 12 {
		containerID = containerID[:12]
	}

	fmt.Printf("Started NATS container %q (ID: %s) with image %q\n", r.containerName, containerID, r.opts.DockerImage)
	fmt.Printf("Connect URL: nats://localhost:%d\n", r.opts.Port)
	fmt.Printf("Monitoring:  http://localhost:8222\n")
	if absStoreDir != "" {
		fmt.Printf("Storage:     %s (volume mounted to /data)\n", absStoreDir)
	}
	return nil
}

// Stop stops and removes the Docker container
func (r *DockerRunner) Stop(ctx context.Context) error {
	dockerBin, err := r.ensureDockerCLI()
	if err != nil {
		return err
	}

	st, err := r.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}

	if st.State == "not created" {
		return fmt.Errorf("docker container %q does not exist", r.containerName)
	}

	if st.Running {
		stopCmd := exec.CommandContext(ctx, dockerBin, "stop", r.containerName)
		var stderr bytes.Buffer
		stopCmd.Stderr = &stderr
		if err := stopCmd.Run(); err != nil {
			return fmt.Errorf("failed to stop container %q: %s", r.containerName, stderr.String())
		}
	}

	// Remove stopped container to maintain clean state
	rmCmd := exec.CommandContext(ctx, dockerBin, "rm", r.containerName)
	var stderr bytes.Buffer
	rmCmd.Stderr = &stderr
	if err := rmCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove container %q: %s", r.containerName, stderr.String())
	}

	fmt.Printf("Stopped and removed container %q\n", r.containerName)
	return nil
}
