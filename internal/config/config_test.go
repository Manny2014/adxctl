package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Call LoadConfig with non-existent config file to ensure it errors
	_, err := LoadConfig("/non/existent/path/config.yaml")
	if err == nil {
		t.Error("expected error when explicit config file does not exist")
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error loading defaults: %v", err)
	}

	if cfg.Nats.Runner != DefaultRunner {
		t.Errorf("expected runner %q, got %q", DefaultRunner, cfg.Nats.Runner)
	}
	if cfg.Nats.Port != DefaultPort {
		t.Errorf("expected port %d, got %d", DefaultPort, cfg.Nats.Port)
	}
	if cfg.Nats.Docker.Image != DefaultDockerImage {
		t.Errorf("expected docker image %q, got %q", DefaultDockerImage, cfg.Nats.Docker.Image)
	}
}

func TestLoadConfig_CustomNatsYaml(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := `
nats:
  runner: "docker"
  port: 5222
  store_dir: "/tmp/custom-nats-data"
  docker:
    image: "nats:alpine"
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	if cfg.Nats.Runner != "docker" {
		t.Errorf("expected runner 'docker', got %q", cfg.Nats.Runner)
	}
	if cfg.Nats.Port != 5222 {
		t.Errorf("expected port 5222, got %d", cfg.Nats.Port)
	}
	if cfg.Nats.StoreDir != "/tmp/custom-nats-data" {
		t.Errorf("expected store_dir '/tmp/custom-nats-data', got %q", cfg.Nats.StoreDir)
	}
	if cfg.Nats.Docker.Image != "nats:alpine" {
		t.Errorf("expected docker image 'nats:alpine', got %q", cfg.Nats.Docker.Image)
	}
}

func TestLoadConfig_ServerFallback(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := `
server:
  runner: "docker"
  port: 6222
  store_dir: "/tmp/server-nats-data"
  docker:
    image: "nats:2.10"
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	if cfg.Nats.Runner != "docker" {
		t.Errorf("expected runner 'docker' via server fallback, got %q", cfg.Nats.Runner)
	}
	if cfg.Nats.Port != 6222 {
		t.Errorf("expected port 6222 via server fallback, got %d", cfg.Nats.Port)
	}
	if cfg.Nats.StoreDir != "/tmp/server-nats-data" {
		t.Errorf("expected store_dir '/tmp/server-nats-data' via server fallback, got %q", cfg.Nats.StoreDir)
	}
	if cfg.Nats.Docker.Image != "nats:2.10" {
		t.Errorf("expected docker image 'nats:2.10' via server fallback, got %q", cfg.Nats.Docker.Image)
	}
}
