package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var natsStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the NATS server",
	Long: `Start the NATS server instance in the background using either local-cli or docker.
Runner mode and settings are read from config.yaml or can be overridden with flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		if err := r.Start(cmd.Context()); err != nil {
			return fmt.Errorf("failed to start NATS server: %w", err)
		}
		return nil
	},
}

func init() {
	natsCmd.AddCommand(natsStartCmd)
}
