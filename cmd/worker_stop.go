package cmd

import (
	"adxctl/internal/registry"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var workerStopCmd = &cobra.Command{
	Use:   "stop [worker-id]",
	Short: "Stop a running worker",
	Long:  `Stop a running worker.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workerID := args[0]

		reg, err := registry.New()
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		processes, err := reg.List()
		if err != nil {
			return fmt.Errorf("failed to list processes: %w", err)
		}

		var proc *registry.Process
		for _, p := range processes {
			if p.ID == workerID {
				proc = &p
				break
			}
		}

		if proc == nil {
			return fmt.Errorf("worker with ID %s not found", workerID)
		}

		process, err := os.FindProcess(proc.PID)
		if err != nil {
			// if process not found, just remove it from registry
			if err := reg.Remove(workerID); err != nil {
				return fmt.Errorf("failed to remove worker from registry: %w", err)
			}
			return fmt.Errorf("failed to find process with PID %d: %w", proc.PID, err)
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to send SIGTERM to process with PID %d: %w", proc.PID, err)
		}

		if err := reg.Remove(workerID); err != nil {
			return fmt.Errorf("failed to remove worker from registry: %w", err)
		}

		fmt.Printf("Worker %s stopped.\n", workerID)
		return nil
	},
}

func init() {
	workerCmd.AddCommand(workerStopCmd)
}
