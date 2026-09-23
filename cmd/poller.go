package cmd

import (
	"github.com/spf13/cobra"
)

// pollerCmd represents the poller command
var pollerCmd = &cobra.Command{
	Use:   "poller",
	Short: "Manage background data pollers",
	Long:  `Manage data pollers (deploy, stop) to fetch events from external systems (like Jira) and publish them to NATS.`,
}

func init() {
	rootCmd.AddCommand(pollerCmd)
}
