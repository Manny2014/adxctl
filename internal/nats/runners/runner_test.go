package nats

import (
	"context"
	"testing"
)

func TestNewRunner_Local(t *testing.T) {
	r, err := New("local-cli", Options{StoreDir: "/tmp/nats-local-store"})
	if err != nil {
		t.Fatalf("unexpected error creating local runner: %v", err)
	}
	lr, ok := r.(*LocalRunner)
	if !ok {
		t.Fatalf("expected *LocalRunner, got %T", r)
	}
	if lr.opts.Port != 4222 {
		t.Errorf("expected default port 4222, got %d", lr.opts.Port)
	}
	if lr.opts.StoreDir != "/tmp/nats-local-store" {
		t.Errorf("expected store_dir '/tmp/nats-local-store', got %q", lr.opts.StoreDir)
	}
}

func TestNewRunner_Docker(t *testing.T) {
	r, err := New("docker", Options{DockerImage: "my-custom/nats:v1", Port: 5222, StoreDir: "/tmp/nats-docker-store"})
	if err != nil {
		t.Fatalf("unexpected error creating docker runner: %v", err)
	}
	dr, ok := r.(*DockerRunner)
	if !ok {
		t.Fatalf("expected *DockerRunner, got %T", r)
	}
	if dr.opts.Port != 5222 {
		t.Errorf("expected port 5222, got %d", dr.opts.Port)
	}
	if dr.opts.DockerImage != "my-custom/nats:v1" {
		t.Errorf("expected image 'my-custom/nats:v1', got %q", dr.opts.DockerImage)
	}
	if dr.opts.StoreDir != "/tmp/nats-docker-store" {
		t.Errorf("expected store_dir '/tmp/nats-docker-store', got %q", dr.opts.StoreDir)
	}
}

func TestNewRunner_Invalid(t *testing.T) {
	_, err := New("invalid-runner", Options{})
	if err == nil {
		t.Error("expected error for invalid runner type")
	}
}

func TestLocalRunner_Status(t *testing.T) {
	lr := NewLocalRunner(Options{Port: 4222, StoreDir: "/tmp/test-nats-store"})
	st, err := lr.Status(context.Background())
	if err != nil {
		t.Fatalf("unexpected error checking status: %v", err)
	}
	if st.Runner != "local-cli" {
		t.Errorf("expected runner 'local-cli', got %q", st.Runner)
	}
	if st.ConnectURL != "nats://localhost:4222" {
		t.Errorf("expected connect URL 'nats://localhost:4222', got %q", st.ConnectURL)
	}
	if st.MonitorURL != "http://localhost:8222" {
		t.Errorf("expected monitor URL 'http://localhost:8222', got %q", st.MonitorURL)
	}
	if st.StoreDir != "/tmp/test-nats-store" {
		t.Errorf("expected store dir '/tmp/test-nats-store', got %q", st.StoreDir)
	}
}
