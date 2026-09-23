package pollers

import (
	"adxctl/pkg/adx"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// GitHubConfig holds configuration for the GitHub poller.
type GitHubConfig struct {
	Repo         string
	GitHubToken  string
	PollInterval time.Duration
}

// RunGitHubPoller starts the GitHub poller.
func RunGitHubPoller(ctx context.Context, cfg GitHubConfig, client *adx.Client) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	lastPollTime := time.Now().UTC().Add(-cfg.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			log.Println("Polling GitHub for new issues...")
			parts := strings.Split(cfg.Repo, "/")
			if len(parts) != 2 {
				log.Printf("Invalid repo format: %s", cfg.Repo)
				continue
			}
			owner, repo := parts[0], parts[1]

			issues, _, err := ghClient.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
				State: "open",
				Since: lastPollTime,
			})
			if err != nil {
				log.Printf("Error fetching issues: %v", err)
				continue
			}

			lastPollTime = time.Now().UTC()

			for _, issue := range issues {
				if issue.IsPullRequest() {
					continue
				}

				issueBytes, err := json.Marshal(issue)
				if err != nil {
					log.Printf("Error marshalling issue: %v", err)
					continue
				}

				task := adx.Task{
					ID:      uuid.New().String(),
					Source:  "github",
					Type:    "issue_created",
					Payload: issueBytes,
				}

				taskBytes, err := json.Marshal(task)
				if err != nil {
					log.Printf("Error marshalling task: %v", err)
					continue
				}

				subject := fmt.Sprintf("tasks.%s.%s", task.Source, task.Type)
				_, err = client.Publish(ctx, subject, taskBytes)
				if err != nil {
					log.Printf("Error publishing task: %v", err)
					continue
				}
				log.Printf("Published task for issue #%d", issue.GetNumber())
			}
		}
	}
}
