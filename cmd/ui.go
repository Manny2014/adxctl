package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var uiPort int
var uiDir string

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Run the adxctl web UI",
	Long:  `Run a web server to host the adxctl single page application (SPA).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := fmt.Sprintf(":%d", uiPort)
		fs := http.FileServer(http.Dir(uiDir))
		http.Handle("/", fs)

		fmt.Printf("Starting UI server on http://localhost%s\n", addr)
		fmt.Printf("Serving files from %s\n", uiDir)
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntVarP(&uiPort, "port", "p", 3000, "Port for the UI server")
	uiCmd.Flags().StringVar(&uiDir, "dir", "./web/build", "Directory to serve UI files from")
}
