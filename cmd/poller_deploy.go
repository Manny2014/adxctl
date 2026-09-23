package cmd

import (
	"adxctl/internal/pollers"
	"adxctl/internal/registry"
	"adxctl/pkg/adx"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	source   string
	interval string
	repo     string
	project  string
	jql      string
)

var pollerDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy and start a new poller",
	Long:  `Deploy and start a new poller for a specific source.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pollInterval, err := time.ParseDuration(interval)
		if err != nil {
			return fmt.Errorf("invalid interval: %w", err)
		}

		natsURL := os.Getenv("NATS_URL")
		if natsURL == "" {
			return fmt.Errorf("NATS_URL environment variable not set")
		}

		natsClient, err := adx.NewClient(adx.NatsConfig{
			URL:  natsURL,
			Name: "adxctl-poller",
		})
		if err != nil {
			return fmt.Errorf("failed to create nats client: %w", err)
		}
		defer natsClient.Close()

		reg, err := registry.New()
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		pollerID := uuid.New().String()
		proc := registry.Process{
			ID:     pollerID,
			Type:   "poller",
			PID:    os.Getpid(),
			Status: "RUNNING",
		}
		if err := reg.Add(proc); err != nil {
			return fmt.Errorf("failed to add poller to registry: %w", err)
		}
		defer reg.Remove(pollerID)

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		fmt.Printf("Poller %s deployed and running with PID %d.\n", pollerID, os.Getpid())

		switch source {
		case "github":
			githubToken := os.Getenv("GITHUB_TOKEN")
			if githubToken == "" {
				return fmt.Errorf("GITHUB_TOKEN environment variable not set")
			}
			cfg := pollers.GitHubConfig{
				Repo:         repo,
				GitHubToken:  githubToken,
				PollInterval: pollInterval,
			}
			return pollers.RunGitHubPoller(ctx, cfg, natsClient)
		case "jira":
			jiraDomain := os.Getenv("JIRA_DOMAIN")
			jiraEmail := os.Getenv("JIRA_EMAIL")
			jiraToken := os.Getenv("JIRA_API_TOKEN")
			if jiraDomain == "" || jiraEmail == "" || jiraToken == "" {
				return fmt.Errorf("JIRA_DOMAIN, JIRA_EMAIL, or JIRA_API_TOKEN environment variable not set")
			}
			cfg := pollers.Config{
				JiraDomain:   jiraDomain,
				JiraEmail:    jiraEmail,
				JiraAPIToken: jiraToken,
				PollInterval: pollInterval,
			}
			return pollers.Run(ctx, cfg, natsClient)
		default:
			return fmt.Errorf("unsupported poller source: %s", source)
		}
	},
}

func init() {
	pollerCmd.AddCommand(pollerDeployCmd)
	pollerDeployCmd.Flags().StringVar(&source, "source", "", "Source type of the poller (e.g. github, jira)")
	pollerDeployCmd.Flags().StringVar(&interval, "interval", "1m", "Polling interval (e.g. 30s, 5m, 1h)")
	pollerDeployCmd.Flags().StringVar(&repo, "repo", "", "GitHub repository in owner/name format")
	pollerDeployCmd.Flags().StringVar(&project, "project", "", "Jira project key")
	pollerDeployCmd.Flags().StringVar(&jql, "jql", "", "JQL query for Jira")
	pollerDeployCmd.MarkFlagRequired("source")
}
