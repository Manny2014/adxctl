package cmd

import (
	"adxctl/internal/registry"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var pollerStopCmd = &cobra.Command{
	Use:   "stop [poller-id]",
	Short: "Stop a running poller",
	Long:  `Stop a running poller.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pollerID := args[0]

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
			if p.ID == pollerID {
				proc = &p
				break
			}
		}

		if proc == nil {
			return fmt.Errorf("poller with ID %s not found", pollerID)
		}

		process, err := os.FindProcess(proc.PID)
		if err != nil {
			// if process not found, just remove it from registry
			if err := reg.Remove(pollerID); err != nil {
				return fmt.Errorf("failed to remove poller from registry: %w", err)
			}
			return fmt.Errorf("failed to find process with PID %d: %w", proc.PID, err)
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to send SIGTERM to process with PID %d: %w", proc.PID, err)
		}

		if err := reg.Remove(pollerID); err != nil {
			return fmt.Errorf("failed to remove poller from registry: %w", err)
		}

		fmt.Printf("Poller %s stopped.\n", pollerID)
		return nil
	},
}

func init() {
	pollerCmd.AddCommand(pollerStopCmd)
}
