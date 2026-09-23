package cmd

import (
	"adxctl/internal/registry"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var workerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running workers",
	Long:  `List all running workers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.New()
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		processes, err := reg.List()
		if err != nil {
			return fmt.Errorf("failed to list processes: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Type", "Subject", "Runtime", "Status"})

		for _, p := range processes {
			if p.Type == "worker" {
				table.Append([]string{p.ID, p.Type, p.Subject, p.Runtime, p.Status})
			}
		}

		table.Render()
		return nil
	},
}

func init() {
	workerCmd.AddCommand(workerListCmd)
}
