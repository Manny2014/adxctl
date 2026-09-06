package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var natsStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running NATS server",
	Long: `Stop the running NATS server instance managed by local-cli or docker.
Runner mode and settings are read from config.yaml or can be overridden with flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		if err := r.Stop(cmd.Context()); err != nil {
			return fmt.Errorf("failed to stop NATS server: %w", err)
		}
		return nil
	},
}

func init() {
	natsCmd.AddCommand(natsStopCmd)
}
