package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var eventbusStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the eventbus server",
	Long:  `Stop the eventbus server instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		if err := r.Stop(cmd.Context()); err != nil {
			return fmt.Errorf("failed to stop eventbus server: %w", err)
		}
		return nil
	},
}

func init() {
	eventbusCmd.AddCommand(eventbusStopCmd)
}
