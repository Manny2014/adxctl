package cmd

import (
	runners "adxctl/internal/nats/runners"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagRunner      string
	flagDockerImage string
	flagPort        int
	flagStoreDir    string
)

// eventbusCmd represents the eventbus command
var eventbusCmd = &cobra.Command{
	Use:   "eventbus",
	Short: "Manage the event bus instance",
	Long:  `Manage the event bus lifecycle (start, stop). Currently supports NATS.`,
}

func init() {
	rootCmd.AddCommand(eventbusCmd)

	eventbusCmd.PersistentFlags().StringVarP(&flagRunner, "runner", "r", "", "server runner: local-cli or docker (default from config.yaml: local-cli)")
	eventbusCmd.PersistentFlags().StringVar(&flagDockerImage, "image", "", "docker image URI when using docker runner (default from config.yaml: nats:latest)")
	eventbusCmd.PersistentFlags().IntVarP(&flagPort, "port", "p", 0, "port for NATS client connections (default from config.yaml: 4222)")
	eventbusCmd.PersistentFlags().StringVar(&flagStoreDir, "store-dir", "", "storage directory for NATS data (default from config.yaml: \"\")")
}

// resolveRunner resolves configuration and creates the appropriate Runner
func resolveRunner(cmd *cobra.Command) (runners.Runner, error) {
	cfg, err := LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	if cfg.EventBus.Type != "nats" {
		return nil, fmt.Errorf("unsupported eventbus type: %s", cfg.EventBus.Type)
	}

	natsCfg := cfg.EventBus.Nats

	runnerType := natsCfg.Runner
	if cmd.Flags().Changed("runner") || flagRunner != "" {
		runnerType = flagRunner
	}

	port := natsCfg.Port
	if cmd.Flags().Changed("port") || flagPort != 0 {
		port = flagPort
	}

	dockerImage := natsCfg.Docker.Image
	if cmd.Flags().Changed("image") || flagDockerImage != "" {
		dockerImage = flagDockerImage
	}

	storeDir := natsCfg.StoreDir
	if cmd.Flags().Changed("store-dir") || flagStoreDir != "" {
		storeDir = flagStoreDir
	}

	opts := runners.Options{
		Port:        port,
		StoreDir:    storeDir,
		DockerImage: dockerImage,
	}

	return runners.New(runnerType, opts)
}
