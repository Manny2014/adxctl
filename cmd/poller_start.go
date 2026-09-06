package cmd

import (
	"adxctl/internal/pollers"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var pollerStartCmd = &cobra.Command{
	Use:   "start [poller-name]",
	Short: "Start a background poller defined in config.yaml",
	Long:  `Start a background poller by its name as defined in the 'pollers' section of your config.yaml.`,
	Args:  cobra.ExactArgs(1), // require exactly one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		pollerName := args[0]

		cfg, err := LoadAppConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		pollerCfg, ok := cfg.Pollers[pollerName]
		if !ok {
			return fmt.Errorf("poller %q not found in configuration", pollerName)
		}

		if !pollerCfg.Enabled {
			return fmt.Errorf("poller %q is disabled in configuration", pollerName)
		}

		switch pollerCfg.Type {
		case "jira":
			jiraCfg := pollers.Config{
				JiraDomain:   pollerCfg.JiraDomain,
				JiraEmail:    pollerCfg.JiraEmail,
				JiraAPIToken: pollerCfg.JiraAPIToken,
				NatsURL:      pollerCfg.NatsURL,
				PollInterval: pollerCfg.Interval,
			}

			// Handle graceful termination
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			return pollers.Run(ctx, jiraCfg)
		default:
			return fmt.Errorf("unsupported poller type %q for poller %q", pollerCfg.Type, pollerName)
		}
	},
}

func init() {
	pollerCmd.AddCommand(pollerStartCmd)
}
