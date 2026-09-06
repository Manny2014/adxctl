package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var natsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the NATS server",
	Long: `Display runtime details of the NATS server instance, including whether it is running,
its connection URL, monitoring endpoint, runner type, PID or container ID, and storage directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := resolveRunner(cmd)
		if err != nil {
			return err
		}

		st, err := r.Status(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to retrieve status: %w", err)
		}

		fmt.Println("NATS Server Status:")
		fmt.Printf("  Runner:      %s\n", st.Runner)
		if st.Running {
			fmt.Printf("  Status:      running\n")
		} else {
			fmt.Printf("  Status:      %s\n", st.State)
		}
		fmt.Printf("  Connect URL: %s\n", st.ConnectURL)
		if st.Running {
			fmt.Printf("  Monitoring:  %s\n", st.MonitorURL)
		}
		if st.PID > 0 {
			fmt.Printf("  PID:         %d\n", st.PID)
		}
		if st.ContainerID != "" {
			fmt.Printf("  Container:   %s (%s)\n", st.ContainerName, st.ContainerID)
		}
		if st.Image != "" && st.Runner == "docker" {
			fmt.Printf("  Image:       %s\n", st.Image)
		}
		if st.StoreDir != "" {
			fmt.Printf("  Storage:     %s\n", st.StoreDir)
		}
		if len(st.Logs) > 0 {
			fmt.Println("\nRecent Logs (last 5 lines):")
			for _, line := range st.Logs {
				fmt.Printf("  %s\n", line)
			}
		}

		return nil
	},
}

func init() {
	natsCmd.AddCommand(natsStatusCmd)
}
