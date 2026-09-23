package cmd

import (
	"adxctl/internal/agent"
	"adxctl/internal/config"
	"adxctl/internal/registry"
	"adxctl/pkg/adx"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var workerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a new worker",
	Long:  `Run a new worker on a specified runtime.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		natsURL := os.Getenv("NATS_URL")
		if natsURL == "" {
			return fmt.Errorf("NATS_URL environment variable not set")
		}

		natsClient, err := adx.NewClient(adx.NatsConfig{
			URL:  natsURL,
			Name: "adxctl-worker",
		})
		if err != nil {
			return fmt.Errorf("failed to create nats client: %w", err)
		}
		defer natsClient.Close()

		var runner agent.Runner
		switch workerRuntime {
		case "local":
			runner, err = agent.NewLocalRunner(config.LocalRunnerConfig{})
			if err != nil {
				return fmt.Errorf("failed to create local runner: %w", err)
			}
		case "docker":
			runner, err = agent.NewDockerRunner(workerImage)
			if err != nil {
				return fmt.Errorf("failed to create docker runner: %w", err)
			}
		default:
			return fmt.Errorf("unsupported runtime: %s", workerRuntime)
		}

		reg, err := registry.New()
		if err != nil {
			return fmt.Errorf("failed to create registry: %w", err)
		}

		workerID := uuid.New().String()
		proc := registry.Process{
			ID:      workerID,
			Type:    "worker",
			PID:     os.Getpid(),
			Subject: workerSubject,
			Runtime: workerRuntime,
			Status:  "RUNNING",
		}
		if err := reg.Add(proc); err != nil {
			return fmt.Errorf("failed to add worker to registry: %w", err)
		}
		defer reg.Remove(workerID)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		fmt.Printf("Worker %s running with PID %d, subscribed to subject: %s\n", workerID, os.Getpid(), workerSubject)

		sub, err := natsClient.NC().QueueSubscribe(workerSubject, "adxctl-workers", func(msg *nats.Msg) {
			var task adx.Task
			if err := json.Unmarshal(msg.Data, &task); err != nil {
				log.Printf("Error unmarshalling task: %v", err)
				return
			}
			if err := runner.Run(ctx, workerSubject, &task); err != nil {
				log.Printf("Error running task %s: %v", task.ID, err)
			}
		})
		if err != nil {
			return fmt.Errorf("failed to subscribe to subject %s: %w", workerSubject, err)
		}
		defer sub.Unsubscribe()

		<-ctx.Done()
		log.Println("Worker shutting down...")

		return nil
	},
}

