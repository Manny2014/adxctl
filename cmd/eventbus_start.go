package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var eventbusStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the eventbus server",
	Long:  `Start the eventbus server instance in the background.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		if err := r.Start(cmd.Context()); err != nil {
			return fmt.Errorf("failed to start eventbus server: %w", err)
		}
		return nil
	},
}

func init() {
	eventbusCmd.AddCommand(eventbusStartCmd)
}
