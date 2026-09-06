package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var apiPort int

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the adxctl API server",
	Long:  `Run the adxctl API server to provide task status and other information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := fmt.Sprintf(":%d", apiPort)
		fmt.Printf("Starting API server on %s...\n", addr)
		// Basic task API endpoint (placeholder)
		http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"id": "task-1", "status": "queued"}, {"id": "task-2", "status": "in_progress"}]`))
		})
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().IntVarP(&apiPort, "port", "p", 8080, "Port for the API server")
}
