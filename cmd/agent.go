package cmd

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage the agent orchestrator",
	Long:  `Manage the agent orchestrator, which consumes events from the event bus and dispatches tasks to runners.`,
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
