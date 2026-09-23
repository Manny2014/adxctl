package adx

import "encoding/json"

// Task represents a single, standardized unit of work that can be processed by a worker.
type Task struct {
	ID       string          `json:"id"`
	Source   string          `json:"source"`
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}
