package cmd

import (
	"adxctl/internal/eventbus"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var runner string

// eventbusStartCmd represents the start command
var eventbusStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the eventbus",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch runner {
		case "docker":
			fmt.Println("Starting event bus with docker runner...")
			if err := eventbus.StartDockerRunner(context.Background()); err != nil {
				return err
			}
			state := eventbus.State{
				RunnerType: "docker",
				Identifier: "adx-eventbus",
			}
			if err := eventbus.SaveState(state); err != nil {
				return err
			}
			fmt.Println("Event bus started successfully.")
		case "local":
			fmt.Println("Starting event bus with local runner...")
			pid, err := eventbus.StartLocalRunner()
			if err != nil {
				return err
			}
			state := eventbus.State{
				RunnerType: "local",
				Identifier: fmt.Sprintf("%d", pid),
			}
			if err := eventbus.SaveState(state); err != nil {
				return err
			}
			fmt.Printf("Event bus started successfully (PID: %d).\n", pid)
		default:
			return fmt.Errorf("invalid runner: %s", runner)
		}
		return nil
	},
}

func init() {
	eventbusCmd.AddCommand(eventbusStartCmd)
	eventbusStartCmd.Flags().StringVarP(&runner, "runner", "r", "docker", "The runner to use. Options are `docker` or `local`.")
}
