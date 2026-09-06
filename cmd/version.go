package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information populated at build time or defaults
var (
	Version   = "0.1.0"
	GitCommit = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of adxctl",
	Long:  `Display version, build commit, build date, and Go runtime information for adxctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("adxctl version %s\n", Version)
		fmt.Printf("  commit:  %s\n", GitCommit)
		fmt.Printf("  built:   %s\n", BuildDate)
		fmt.Printf("  go:      %s\n", runtime.Version())
		fmt.Printf("  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
