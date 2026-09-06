package cmd

import (
	"adxctl/internal/agent"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the agent orchestrator",
	Long:  `Start the agent orchestrator to consume events and dispatch tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Orchestrator setup and execution will go here.
		// For now, this is a placeholder.
		orchestrator, err := agent.NewOrchestrator()
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return orchestrator.Run(ctx)
	},
}

func init() {
	agentCmd.AddCommand(agentStartCmd)
}
