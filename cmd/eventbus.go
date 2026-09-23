package cmd

import (
	"github.com/spf13/cobra"
)

// eventbusCmd represents the eventbus command
var eventbusCmd = &cobra.Command{
	Use:   "eventbus",
	Short: "Manage the local development event bus",
	Long:  `Start and stop a local NATS event bus for development purposes.`,
}

func init() {
	rootCmd.AddCommand(eventbusCmd)
}
