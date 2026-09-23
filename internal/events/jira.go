package events

import "time"

// Config holds the server connection and bucket setup details.
type Config struct {
	URL string        // e.g.,  https://<your-domain>.atlassian.net/rest/api/3/
}

type JiraEvent struct {
	Config Config  `json:"config"`
	Domain string `json:"domain"`
}

func (je *JiraEvent) EventType() string {
    return "jira"
}

func (je *JiraEvent) OccurredAt() time.Time {
    return time.Now()
}