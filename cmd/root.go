package cmd

import (
	"fmt"
	"os"

	"adxctl/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// LoadAppConfig loads the configuration based on the --config flag or default search paths
func LoadAppConfig() (*config.Config, error) {
	return config.LoadConfig(cfgFile)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "adxctl",
	Short:         "adxctl is a CLI tool for managing ADX workflows",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `adxctl is a fast and flexible command-line interface built with Go and Cobra.
It provides commands and utilities for managing ADX resources, automation, and tasks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is ./config.yaml or $HOME/.adxctl/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
