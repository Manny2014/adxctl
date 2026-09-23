package events

import "encoding/json"

// Poller represents a configurable, long-running process for fetching data.
type Poller struct {
	ID       string          `json:"id"`
	Source   string          `json:"source"`
	Interval string          `json:"interval"`
	Config   json.RawMessage `json:"config"`
	Status   string          `json:"status"`
}

// Worker represents a process that executes tasks.
type Worker struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Runtime string `json:"runtime"`
	Status  string `json:"status"`
}
