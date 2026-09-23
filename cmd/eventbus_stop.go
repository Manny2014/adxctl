package cmd

import (
	"adxctl/internal/eventbus"
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// eventbusStopCmd represents the stop command
var eventbusStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the eventbus",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := eventbus.LoadState()
		if err != nil {
			return fmt.Errorf("no adxctl-managed event bus is currently running")
		}

		switch state.RunnerType {
		case "docker":
			fmt.Println("Stopping event bus (docker runner)...")
			if err := eventbus.StopDockerRunner(context.Background()); err != nil {
				return err
			}
			if err := eventbus.ClearState(); err != nil {
				return err
			}
			fmt.Println("Event bus stopped successfully.")
		case "local":
			fmt.Printf("Stopping event bus (local runner, PID: %s)...\n", state.Identifier)
			pid, err := strconv.Atoi(state.Identifier)
			if err != nil {
				return fmt.Errorf("invalid PID in state file: %s", state.Identifier)
			}
			if err := eventbus.StopLocalRunner(pid); err != nil {
				return err
			}
			if err := eventbus.ClearState(); err != nil {
				return err
			}
			fmt.Println("Event bus stopped successfully.")
		default:
			return fmt.Errorf("unknown runner type in state file: %s", state.RunnerType)
		}
		return nil
	},
}

func init() {
	eventbusCmd.AddCommand(eventbusStopCmd)
}
