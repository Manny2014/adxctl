package cmd

import (
	"fmt"

	nats "adxctl/internal/nats/runners"

	"github.com/spf13/cobra"
)

var (
	flagRunner      string
	flagDockerImage string
	flagPort        int
	flagStoreDir    string
)

// natsCmd represents the nats command
var natsCmd = &cobra.Command{
	Use:   "nats",
	Short: "Manage the NATS server instance",
	Long: `Manage the NATS server lifecycle (start, stop).
Supports running locally via local-cli (nats-server binary) or containerized via Docker.
Runner behavior can be configured in config.yaml under 'nats' or overridden via CLI flags.`,
}

func init() {
	rootCmd.AddCommand(natsCmd)

	natsCmd.PersistentFlags().StringVarP(&flagRunner, "runner", "r", "", "server runner: local-cli or docker (default from config.yaml: local-cli)")
	natsCmd.PersistentFlags().StringVar(&flagDockerImage, "image", "", "docker image URI when using docker runner (default from config.yaml: nats:latest)")
	natsCmd.PersistentFlags().IntVarP(&flagPort, "port", "p", 0, "port for NATS client connections (default from config.yaml: 4222)")
	natsCmd.PersistentFlags().StringVar(&flagStoreDir, "store-dir", "", "storage directory for NATS data (default from config.yaml: \"\")")
}

// resolveRunner resolves configuration and creates the appropriate Runner
func resolveRunner(cmd *cobra.Command) (nats.Runner, error) {
	cfg, err := LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	runnerType := cfg.Nats.Runner
	if cmd.Flags().Changed("runner") || flagRunner != "" {
		runnerType = flagRunner
	}

	port := cfg.Nats.Port
	if cmd.Flags().Changed("port") || flagPort != 0 {
		port = flagPort
	}

	dockerImage := cfg.Nats.Docker.Image
	if cmd.Flags().Changed("image") || flagDockerImage != "" {
		dockerImage = flagDockerImage
	}

	storeDir := cfg.Nats.StoreDir
	if cmd.Flags().Changed("store-dir") || flagStoreDir != "" {
		storeDir = flagStoreDir
	}

	opts := nats.Options{
		Port:        port,
		StoreDir:    storeDir,
		DockerImage: dockerImage,
	}

	return nats.New(runnerType, opts)
}
