package nats

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Config defines connection and operational settings for the Client.
type NatsConfig struct {
	URL            string        // e.g. "nats://localhost:4222" or nats.DefaultURL
	ConnectTimeout time.Duration // Timeout for establishing connection (default: 5s)
	MaxReconnects  int           // Max reconnection attempts (-1 for unlimited, default: -1)
	ReconnectWait  time.Duration // Wait time between reconnects (default: 2s)
	Name           string        // Client connection name for NATS monitoring
}

// Client encapsulates NATS Core, JetStream, and helper features.
type Client struct {
	nc *nats.Conn
	js jetstream.JetStream
	mu sync.RWMutex
}

// NewClient initializes and returns a ready-to-use NATS JetStream client.
func NewClient(cfg NatsConfig) (*Client, error) {
	if cfg.URL == "" {
		cfg.URL = nats.DefaultURL
	}
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 5 * time.Second
	}
	if cfg.ReconnectWait == 0 {
		cfg.ReconnectWait = 2 * time.Second
	}
	if cfg.MaxReconnects == 0 {
		cfg.MaxReconnects = -1 // Default to infinite retries
	}

	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.Timeout(cfg.ConnectTimeout),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("[NATS] Client disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[NATS] Client reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Println("[NATS] Connection closed")
		}),
	}

	// 1. Establish NATS Connection
	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %w", err)
	}

	// 2. Initialize JetStream Context
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream init failed: %w", err)
	}

	return &Client{
		nc: nc,
		js: js,
	}, nil
}

// JS returns the underlying jetstream.JetStream interface.
func (c *Client) JS() jetstream.JetStream {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.js
}

// NC returns the raw nats.Conn for standard core Pub/Sub operations.
func (c *Client) NC() *nats.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nc
}

// EnsureStream idempotently creates or updates a JetStream stream.
func (c *Client) EnsureStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("stream name cannot be empty")
	}
	return c.js.CreateOrUpdateStream(ctx, cfg)
}

// EnsureKV idempotently creates or updates a Key-Value bucket.
func (c *Client) EnsureKV(ctx context.Context, cfg jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}
	return c.js.CreateOrUpdateKeyValue(ctx, cfg)
}

// Publish convenient wrapper to publish messages to a JetStream subject.
func (c *Client) Publish(ctx context.Context, subject string, data []byte) (*jetstream.PubAck, error) {
	return c.js.Publish(ctx, subject, data)
}

// Close gracefully drains pending messages and closes the underlying connection.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nc != nil && !c.nc.IsClosed() {
		if err := c.nc.Drain(); err != nil {
			log.Printf("[NATS] Error during drain: %v", err)
		}
	}
}


/* Example Usage...
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"your_module/nats" // Update with your Go module path
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Initialize the Client
	client, err := nats.NewClient(nats.Config{
		URL:            "nats://127.0.0.1:4222",
		Name:           "jira-integration-service",
		ConnectTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to initialize NATS client: %v", err)
	}
	defer client.Close()

	// 2. Ensure Stream Exists
	_, err = client.EnsureStream(ctx, jetstream.StreamConfig{
		Name:     "JIRA_EVENTS",
		Subjects: []string{"jira.>"},
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// 3. Publish a Message
	ack, err := client.Publish(ctx, "jira.issue_created.PROJ", []byte(`{"id": "10001"}`))
	if err != nil {
		log.Fatalf("Publish error: %v", err)
	}
	fmt.Printf("Message published successfully! Sequence: %d\n", ack.Sequence)

	// 4. Ensure KV Bucket Exists
	kv, err := client.EnsureKV(ctx, jetstream.KeyValueConfig{
		Bucket: "poller_watermarks",
		TTL:    24 * time.Hour,
	})
	if err != nil {
		log.Fatalf("KV error: %v", err)
	}

	// Put/Get from KV
	kv.Put(ctx, "last_run", []byte(time.Now().Format(time.RFC3339)))
}
*/