package cmd

import (
	"github.com/spf13/cobra"
)

var (
	workerSubject string
	workerRuntime string
	workerImage   string
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Manage workers",
	Long:  `Manage workers that process tasks from the event bus.`,
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.AddCommand(workerRunCmd)
	workerRunCmd.Flags().StringVar(&workerSubject, "subject", "", "NATS subject to subscribe to")
	workerRunCmd.Flags().StringVar(&workerRuntime, "runtime", "local", "Worker runtime (local, docker)")
	workerRunCmd.Flags().StringVar(&workerImage, "image", "adxctl-worker", "Docker image to use for the worker")
	workerRunCmd.MarkFlagRequired("subject")
}
