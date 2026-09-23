package pollers

import (
	"adxctl/pkg/adx"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// Config holds environment and service runtime parameters.
type Config struct {
	JiraDomain   string
	JiraEmail    string
	JiraAPIToken string
	PollInterval time.Duration
}

// JiraSearchResponse models the relevant fields from Jira REST API v3 /search response.
type JiraSearchResponse struct {
	Issues []JiraIssue `json:"issues"`
}

type JiraIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary   string `json:"summary"`
		Created   string `json:"created"`
		Updated   string `json:"updated"`
		IssueType struct {
			ID      string `json:"id"`
			Name    string `json:"name"`    // e.g., "Sub-task", "Epic", "Task", "Bug", "Story"
			Subtask bool   `json:"subtask"` // true if it is a sub-task
		} `json:"issuetype"`
		Parent *struct {
			ID  string `json:"id"`
			Key string `json:"key"` // Key of the parent issue if this is a sub-task
		} `json:"parent,omitempty"`
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
	} `json:"fields"`
}

const (
	StreamName      = "JIRA_EVENTS"
	WatermarkBucket = "jira_poller_state"
	WatermarkKey    = "last_poll_timestamp"
)

type JiraSearchRequest struct {
	JQL        string   `json:"jql"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields,omitempty"`
}

// Run connects to NATS, configures JetStream streams/buckets, and starts the polling loop.
func Run(ctx context.Context, cfg Config, client *adx.Client) error {
	// FIX 1: Bind stream to tasks.jira.> to capture published subjects "tasks.jira.<eventType>"
	log.Printf("[INIT] Ensuring JetStream stream %q exists with subjects %v...", StreamName, []string{"tasks.jira.>"})
	_, err := client.EnsureStream(ctx, jetstream.StreamConfig{
		Name:        StreamName,
		Description: "Polled Jira events",
		Subjects:    []string{"tasks.jira.>"},
	})

	if err != nil {
		return fmt.Errorf("failed to setup JetStream stream: %w", err)
	}
	log.Printf("[INIT] JetStream stream %q is verified and ready.", StreamName)

	log.Printf("[INIT] Ensuring state Key-Value bucket %q exists...", WatermarkBucket)

	kv, err := client.EnsureKV(ctx, jetstream.KeyValueConfig{
		Bucket:      WatermarkBucket,
		Description: "Stores the last poll timestamp for state recovery",
	})

	if err != nil {
		return fmt.Errorf("failed to setup KV bucket: %w", err)
	}
	log.Printf("[INIT] Key-Value bucket %q is verified and ready.", WatermarkBucket)

	log.Printf("[INIT] Jira Poller successfully started for domain %q. Polling interval: %s.", cfg.JiraDomain, cfg.PollInterval)

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	// FIX 2: Trigger initial run with a distinct context deadline
	log.Println("[POLL] Triggering initial polling execution cycle...")
	runPoll(ctx, cfg, client, kv)

	for {
		select {
		case <-ctx.Done():
			log.Println("[SHUTDOWN] Shutting down Jira Poller gracefully...")
			return nil
		case <-ticker.C:
			log.Println("[POLL] Triggering scheduled polling execution cycle...")
			// FIX 2: Instantiate a fresh context per ticker interval
			runPoll(ctx, cfg, client, kv)
		}
	}
}

// Helper wrapper to ensure fresh per-execution contexts
func runPoll(parentCtx context.Context, cfg Config, client *adx.Client, kv jetstream.KeyValue) {
	pollCtx, cancel := context.WithTimeout(parentCtx, 15*time.Second)
	defer cancel()
	pollAndPublish(pollCtx, cfg, client, kv)
}

func pollAndPublish(ctx context.Context, cfg Config, client *adx.Client, kv jetstream.KeyValue) {
	lastPollTime := getLastPollTime(ctx, kv)
	currentPollTime := time.Now().UTC()

	jqlTimeStr := lastPollTime.Format("2006-01-02 15:04")
	jql := fmt.Sprintf(`updated >= "%s" ORDER BY updated ASC`, jqlTimeStr)

	log.Printf("[POLL] Executing Jira sync. Querying JQL: %s", jql)

	issues, err := fetchJiraIssues(ctx, cfg, jql)
	if err != nil {
		log.Printf("[ERROR] Jira sync query failed: %v", err)
		return
	}

	log.Printf("[POLL] Retrieved %d updated ticket(s) modified since last watermark (%s).", len(issues), jqlTimeStr)

	if len(issues) == 0 {
		setLastPollTime(ctx, kv, currentPollTime)
		return
	}

	latestIssueUpdatedTime := lastPollTime

	for i, issue := range issues {
		projectKey := issue.Fields.Project.Key
		if projectKey == "" {
			projectKey = "UNKNOWN"
		}

		var eventType string
		createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", issue.Fields.Created)
		if err != nil {
			log.Printf("[WARN] [%d/%d] Could not parse 'created' timestamp for issue %s. Defaulting to 'issue_updated'. Error: %v", i+1, len(issues), issue.Key, err)
			eventType = "issue_updated"
		} else {
			if createdTime.After(lastPollTime) {
				eventType = "issue_created"
			} else {
				eventType = "issue_updated"
			}
		}

		// Published subject: tasks.jira.issue_created / tasks.jira.issue_updated
		subject := fmt.Sprintf("tasks.jira.%s", eventType)

		issueBytes, err := json.Marshal(issue)
		if err != nil {
			log.Printf("[ERROR] [%d/%d] Failed to marshal issue for %s %s: %v",
				i+1, len(issues), issue.Fields.IssueType.Name, issue.Key, err)
			continue
		}

		task := adx.Task{
			ID:      uuid.New().String(),
			Source:  "jira",
			Type:    eventType,
			Payload: issueBytes,
		}

		taskBytes, err := json.Marshal(task)
		if err != nil {
			log.Printf("[ERROR] [%d/%d] Failed to marshal task for %s %s: %v",
				i+1, len(issues), issue.Fields.IssueType.Name, issue.Key, err)
			continue
		}

		log.Printf("[POLL] [%d/%d] Publishing '%s' event for %s %s to NATS subject %q...",
			i+1, len(issues), eventType, issue.Fields.IssueType.Name, issue.Key, subject)

		ack, err := client.Publish(ctx, subject, taskBytes)
		if err != nil {
			log.Printf("[ERROR] [%d/%d] JetStream publish failed for %s %s: %v",
				i+1, len(issues), issue.Fields.IssueType.Name, issue.Key, err)
			continue
		}

		log.Printf("[POLL] [%d/%d] Successfully published %s %s [%s] to NATS (Stream Seq: %d).",
			i+1, len(issues), issue.Fields.IssueType.Name, issue.Key, issue.Fields.Summary, ack.Sequence)

		if parsedTime, err := time.Parse("2006-01-02T15:04:05.000-0700", issue.Fields.Updated); err == nil {
			if parsedTime.After(latestIssueUpdatedTime) {
				latestIssueUpdatedTime = parsedTime
			}
		}
	}

	if latestIssueUpdatedTime.After(lastPollTime) {
		setLastPollTime(ctx, kv, latestIssueUpdatedTime)
	} else {
		setLastPollTime(ctx, kv, currentPollTime)
	}
}

func fetchJiraIssues(ctx context.Context, cfg Config, jql string) ([]JiraIssue, error) {
	reqURL := fmt.Sprintf("https://%s/rest/api/3/search/jql", cfg.JiraDomain)

	payload := JiraSearchRequest{
		JQL:        jql,
		MaxResults: 50,
		Fields:     []string{"*all"},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	log.Printf("[HTTP] Preparing POST request to JIRA API search endpoint: %s", reqURL)
	log.Printf("[HTTP] JQL Payload: %s", string(jsonBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(cfg.JiraEmail, cfg.JiraAPIToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[HTTP] Dispatching request to Jira...")
	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)
	log.Printf("[HTTP] Response received in %s. Status Code: %d", duration, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var searchResp JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	log.Printf("[HTTP] Successfully parsed %d issue(s) from the search response.", len(searchResp.Issues))
	return searchResp.Issues, nil
}

func getLastPollTime(ctx context.Context, kv jetstream.KeyValue) time.Time {
	log.Printf("[STATE] Fetching last sync watermark from Key-Value bucket using key %q...", WatermarkKey)
	entry, err := kv.Get(ctx, WatermarkKey)
	if err == nil {
		if parsed, err := time.Parse(time.RFC3339, string(entry.Value())); err == nil {
			log.Printf("[STATE] Successfully retrieved and parsed watermark: %s", parsed.Format(time.RFC3339))
			return parsed
		}
	}

	fallback := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	log.Printf("[STATE] No previous sync state found (or failed to parse, err: %v). Using fallback watermark (beginning of 2026): %s", err, fallback.Format(time.RFC3339))
	return fallback
}

func setLastPollTime(ctx context.Context, kv jetstream.KeyValue, t time.Time) {
	val := t.Format(time.RFC3339)
	log.Printf("[STATE] Saving new sync watermark to Key-Value bucket: %s...", val)
	_, err := kv.Put(ctx, WatermarkKey, []byte(val))
	if err != nil {
		log.Printf("[WARN] Failed to persist new sync state: %v", err)
	} else {
		log.Printf("[STATE] Successfully persisted watermark state: %s", val)
	}
}