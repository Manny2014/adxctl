package cmd

import (
	"adxctl/internal/registry"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var pollerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running pollers",
	Long:  `List all running pollers.`,
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
		table.Header([]string{"ID", "Type", "Status"})

		for _, p := range processes {
			if p.Type == "poller" {
				table.Append([]string{p.ID, p.Type, p.Status})
			}
		}

		table.Render()
		return nil
	},
}

func init() {
	pollerCmd.AddCommand(pollerListCmd)
}
