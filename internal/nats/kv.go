package nats

import (
	"adxctl/pkg/adx"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// KVConfig holds the server connection and bucket setup details.
type KVConfig struct {
	adx.NatsConfig
	Bucket string        // Cache bucket name
	TTL    time.Duration // Default TTL for entries in this bucket (0 = no expiry)
}

// KVClient wraps NATS JetStream Key-Value store with typical cache APIs.
type KVClient struct {
	client *adx.Client
	kv     jetstream.KeyValue
}

// NewKVClient connects to NATS, initializes JetStream, and creates/binds to the bucket.
func NewKVClient(ctx context.Context, cfg KVConfig) (*KVClient, error) {
	if cfg.Bucket == "" {
		return nil, errors.New("bucket name is required")
	}

	client, err := adx.NewClient(cfg.NatsConfig)
	if err != nil {
		return nil, fmt.Errorf("nats client init failed: %w", err)
	}

	kv, err := client.EnsureKV(ctx, jetstream.KeyValueConfig{
		Bucket:      cfg.Bucket,
		Description: "Application Cache Store",
		TTL:         cfg.TTL,
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("kv bucket setup failed: %w", err)
	}

	return &KVClient{
		client: client,
		kv:     kv,
	}, nil
}

// Set stores a key-value pair in the cache.
func (c *KVClient) Set(ctx context.Context, key string, val []byte) error {
	_, err := c.kv.Put(ctx, key, val)
	if err != nil {
		return fmt.Errorf("cache set failed for key '%s': %w", key, err)
	}
	return nil
}

// Get retrieves a value by key. Returns nil if the key does not exist.
func (c *KVClient) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := c.kv.Get(ctx, key)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, nil // Return nil slice on key miss
		}
		return nil, fmt.Errorf("cache get failed for key '%s': %w", key, err)
	}
	return entry.Value(), nil
}

// Delete removes a key from the cache.
func (c *KVClient) Delete(ctx context.Context, key string) error {
	err := c.kv.Delete(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return fmt.Errorf("cache delete failed for key '%s': %w", key, err)
	}
	return nil
}

// GetOrSet attempts to get a key; if absent, populates it using factory and returns the value.
func (c *KVClient) GetOrSet(ctx context.Context, key string, factory func() ([]byte, error)) ([]byte, error) {
	val, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if val != nil {
		return val, nil
	}

	// Cache miss: generate value
	newVal, err := factory()
	if err != nil {
		return nil, fmt.Errorf("factory generation failed: %w", err)
	}

	if err := c.Set(ctx, key, newVal); err != nil {
		return nil, err
	}

	return newVal, nil
}

// Close gracefully closes the underlying NATS connection.
func (c *KVClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// EXAMPLE:
// ---------------------------------------------------------
// Example Usage
// ---------------------------------------------------------
// func main() {
// 	ctx := context.Background()

// 	client, err := NewKVClient(ctx, KVConfig{
// 		NatsConfig: NatsConfig{
// 			URL:            "nats://127.0.0.1:4222",
// 			ConnectTimeout: 3 * time.Second,
// 		},
// 		Bucket:         "app_cache",
// 		TTL:            30 * time.Minute,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to initialize KVClient: %v", err)
// 	}
// 	defer client.Close()

// 	key := "session:99"

// 	// 1. Set key
// 	err = client.Set(ctx, key, []byte("active"))
// 	if err != nil {
// 		log.Fatalf("Set error: %v", err)
// 	}

// 	// 2. Get key
// 	val, err := client.Get(ctx, key)
// 	if err != nil {
// 		log.Fatalf("Get error: %v", err)
// 	}
// 	fmt.Printf("Fetched Value: %s\n", string(val))

// 	// 3. GetOrSet pattern
// 	cachedVal, err := client.GetOrSet(ctx, "session:100", func() ([]byte, error) {
// 		// Simulating DB hit
// 		return []byte("created_from_factory"), nil
// 	})
// 	if err != nil {
// 		log.Fatalf("GetOrSet error: %v", err)
// 	}
// 	fmt.Printf("GetOrSet Value: %s\n", string(cachedVal))

// 	// 4. Delete key
// 	if err := client.Delete(ctx, key); err != nil {
// 		log.Fatalf("Delete error: %v", err)
// 	}
// }
