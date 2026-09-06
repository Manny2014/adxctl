package nats

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// LocalRunner manages a nats-server instance running directly on the host machine
type LocalRunner struct {
	opts Options
}

// NewLocalRunner constructs a LocalRunner
func NewLocalRunner(opts Options) *LocalRunner {
	return &LocalRunner{opts: opts}
}

func (r *LocalRunner) getRunDir(create bool) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	dir := filepath.Join(home, ".adxctl")
	if create {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return dir, nil
}

func (r *LocalRunner) pidFilePath(create bool) (string, error) {
	dir, err := r.getRunDir(create)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nats.pid"), nil
}

func (r *LocalRunner) logFilePath(create bool) (string, error) {
	dir, err := r.getRunDir(create)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nats.log"), nil
}

func (r *LocalRunner) readPID() (int, error) {
	pidPath, err := r.pidFilePath(false)
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid pid in %s: %w", pidPath, err)
	}
	return pid, nil
}

// tailFile returns the last n lines from a file. Returns empty slice if unavailable.
func tailFile(path string, n int) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines
}

// isProcessRunning checks if a process with given PID is still active
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix systems, FindProcess always succeeds. Send signal 0 to verify existence.
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// Status checks if local nats-server is running
func (r *LocalRunner) Status(ctx context.Context) (*StatusInfo, error) {
	info := &StatusInfo{
		Runner:     "local-cli",
		Port:       r.opts.Port,
		ConnectURL: fmt.Sprintf("nats://localhost:%d", r.opts.Port),
		MonitorURL: "http://localhost:8222",
		StoreDir:   r.opts.StoreDir,
		State:      "stopped",
		Running:    false,
	}

	pid, err := r.readPID()
	if err != nil {
		if os.IsNotExist(err) {
			return info, nil
		}
		return nil, err
	}

	if isProcessRunning(pid) {
		info.Running = true
		info.State = "running"
		info.PID = pid
		// Attach last 5 log lines when running
		if logPath, err := r.logFilePath(false); err == nil {
			info.Logs = tailFile(logPath, 5)
		}
		return info, nil
	}

	// PID file exists but process is gone: clean up stale PID file
	pidPath, _ := r.pidFilePath(false)
	_ = os.Remove(pidPath)
	info.State = "stopped (stale PID removed)"
	return info, nil
}

// Start launches nats-server in the background
func (r *LocalRunner) Start(ctx context.Context) error {
	binPath, err := exec.LookPath("nats-server")
	if err != nil {
		return fmt.Errorf("nats-server executable not found in PATH. Install via 'brew install nats-server' or download from https://github.com/nats-io/nats-server/releases")
	}

	st, err := r.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to check status: %w", err)
	}
	if st.Running {
		return fmt.Errorf("nats-server is already running (PID %d)", st.PID)
	}

	logPath, err := r.logFilePath(true)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", logPath, err)
	}
	defer logFile.Close()

	portStr := strconv.Itoa(r.opts.Port)
	args := []string{"-js","-p", portStr, "-m", "8222"}

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
		args = append(args, "-sd", absStoreDir)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	// Detach process so it continues running after adxctl exits
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start nats-server: %w", err)
	}

	// Wait a moment to ensure process didn't instantly terminate
	time.Sleep(100 * time.Millisecond)
	if !isProcessRunning(cmd.Process.Pid) {
		return fmt.Errorf("nats-server failed to start. Check logs at %s", logPath)
	}

	pidPath, err := r.pidFilePath(true)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file %s: %w", pidPath, err)
	}

	fmt.Printf("Started local NATS server (PID: %d)\n", cmd.Process.Pid)
	fmt.Printf("Connect URL: nats://localhost:%d\n", r.opts.Port)
	fmt.Printf("Monitoring:  http://localhost:8222\n")
	if absStoreDir != "" {
		fmt.Printf("Storage:     %s\n", absStoreDir)
	}
	fmt.Printf("Logs:        %s\n", logPath)
	return nil
}

// Stop stops the running nats-server process
func (r *LocalRunner) Stop(ctx context.Context) error {
	pid, err := r.readPID()
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("nats-server is not running (no PID file found)")
		}
		return err
	}

	pidPath, _ := r.pidFilePath(false)
	defer os.Remove(pidPath)

	if !isProcessRunning(pid) {
		return fmt.Errorf("process with PID %d is not running", pid)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to process %d: %w", pid, err)
	}

	// Wait up to 5 seconds for process termination
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !isProcessRunning(pid) {
			fmt.Printf("Stopped local nats-server (PID: %d)\n", pid)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Force kill if not terminated
	_ = process.Signal(syscall.SIGKILL)
	fmt.Printf("Forcefully terminated local nats-server (PID: %d)\n", pid)
	return nil
}
