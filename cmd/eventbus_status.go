package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var eventbusStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the eventbus server",
	Long:  `Get the status of the eventbus server instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		status, err := r.Status(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get eventbus status: %w", err)
		}
		fmt.Println(status)
		return nil
	},
}

func init() {
	eventbusCmd.AddCommand(eventbusStatusCmd)
}
