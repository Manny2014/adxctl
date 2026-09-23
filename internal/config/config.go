package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// EventBusConfig holds config for the event bus.
type EventBusConfig struct {
	Type string     `mapstructure:"type"`
	Nats NatsConfig `mapstructure:"nats"`
}

// Config holds the root configuration for adxctl
type Config struct {
	EventBus EventBusConfig          `mapstructure:"eventbus"`
	Pollers  map[string]PollerConfig `mapstructure:"pollers"`
	Agent    AgentConfig             `mapstructure:"agent"`

	// Deprecated: use EventBus section
	Nats   NatsConfig `mapstructure:"nats"`
	Server NatsConfig `mapstructure:"server"`
}

// PollerConfig defines the configuration for a single poller instance
type PollerConfig struct {
	Type     string        `mapstructure:"type"`
	Enabled  bool          `mapstructure:"enabled"`
	Interval time.Duration `mapstructure:"interval"`
	NatsURL  string        `mapstructure:"nats_url"`

	// Jira-specific settings (used if type is 'jira')
	JiraDomain   string `mapstructure:"domain,omitempty"`
	JiraEmail    string `mapstructure:"email,omitempty"`
	JiraAPIToken string `mapstructure:"api_token,omitempty"`
}

// NatsConfig holds configuration for the NATS server
type NatsConfig struct {
	Runner   string       `mapstructure:"runner" yaml:"runner"`
	Port     int          `mapstructure:"port" yaml:"port"`
	StoreDir string       `mapstructure:"store_dir" yaml:"store_dir"`
	Docker   DockerConfig `mapstructure:"docker" yaml:"docker"`
}

// DockerConfig holds Docker runner specific configuration
type DockerConfig struct {
	Image string `mapstructure:"image" yaml:"image"`
}

// AgentConfig holds configuration for the agent system
type AgentConfig struct {
	Runners    AgentRunnersConfig `mapstructure:"runners"`
	AIProvider AIProviderConfig   `mapstructure:"ai_provider"`
}

// AgentRunnersConfig holds configuration for the different task runners
type AgentRunnersConfig struct {
	Default    string               `mapstructure:"default"`
	Local      LocalRunnerConfig    `mapstructure:"local"`
	Docker     DockerRunnerConfig   `mapstructure:"docker"`
	Kubernetes KubernetesRunnerConfig `mapstructure:"kubernetes"`
}

// LocalRunnerConfig holds configuration for the local runner.
type LocalRunnerConfig struct {
	Type   string `mapstructure:"type"`
	Script string `mapstructure:"script,omitempty"` // Used when type is "local-script"
}

// DockerRunnerConfig is a placeholder for docker runner settings
type DockerRunnerConfig struct{}

// KubernetesRunnerConfig is a placeholder for kubernetes runner settings
type KubernetesRunnerConfig struct{}

// AIProviderConfig holds configuration for the AI model provider
type AIProviderConfig struct {
	Type   string         `mapstructure:"type"`
	Google GoogleAIConfig `mapstructure:"google"`
}

// GoogleAIConfig holds configuration for the Google AI provider
type GoogleAIConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

const (
	RunnerLocalCLI     = "local-cli"
	RunnerDocker       = "docker"
	DefaultRunner      = RunnerLocalCLI
	DefaultPort        = 4222
	DefaultDockerImage = "nats:latest"
	DefaultStoreDir    = ""
)

// SetDefaults configures default values for NATS settings
func SetDefaults(v *viper.Viper) {
	// New eventbus defaults
	v.SetDefault("eventbus.type", "nats")
	v.SetDefault("eventbus.nats.runner", DefaultRunner)
	v.SetDefault("eventbus.nats.port", DefaultPort)
	v.SetDefault("eventbus.nats.store_dir", DefaultStoreDir)
	v.SetDefault("eventbus.nats.docker.image", DefaultDockerImage)

	// Agent defaults
	v.SetDefault("agent.runners.default", "local")
	v.SetDefault("agent.ai_provider.type", "google")
	v.SetDefault("agent.ai_provider.google.model", "gemini-pro")

	// Old defaults for backward compatibility
	v.SetDefault("nats.runner", DefaultRunner)
	v.SetDefault("nats.port", DefaultPort)
	v.SetDefault("nats.store_dir", DefaultStoreDir)
	v.SetDefault("nats.docker.image", DefaultDockerImage)

	// Fallback keys for server.*
	v.SetDefault("server.runner", DefaultRunner)
	v.SetDefault("server.port", DefaultPort)
	v.SetDefault("server.store_dir", DefaultStoreDir)
	v.SetDefault("server.docker.image", DefaultDockerImage)
}

// LoadConfig loads the configuration from file, environment, or defaults
func LoadConfig(cfgFile string) (*Config, error) {
	v := viper.New()
	SetDefaults(v)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		// Look for config.yaml in current directory first
		v.AddConfigPath(".")
		// Fallback to ~/.adxctl/
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".adxctl"))
		}
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.SetEnvPrefix("ADXCTL")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && cfgFile != "" {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Handle backward compatibility for nats/server config
	if v.InConfig("server") && !v.InConfig("nats") {
		if cfg.Server.Runner != "" {
			cfg.Nats.Runner = cfg.Server.Runner
		}
		if cfg.Server.Port != 0 {
			cfg.Nats.Port = cfg.Server.Port
		}
		if cfg.Server.StoreDir != "" {
			cfg.Nats.StoreDir = cfg.Server.StoreDir
		}
		if cfg.Server.Docker.Image != "" {
			cfg.Nats.Docker.Image = cfg.Server.Docker.Image
		}
	}

	if (v.InConfig("nats") || v.InConfig("server")) && !v.InConfig("eventbus") {
		cfg.EventBus.Type = "nats"
		cfg.EventBus.Nats = cfg.Nats
	}

	// Expand environment variables in poller configurations
	for name, pollerCfg := range cfg.Pollers {
		expandedPollerCfg := expandPollerEnvVars(pollerCfg)
		cfg.Pollers[name] = expandedPollerCfg
	}

	// Expand environment variables in agent configuration
	cfg.Agent.AIProvider.Google.APIKey = os.ExpandEnv(cfg.Agent.AIProvider.Google.APIKey)

	return &cfg, nil
}

// expandPollerEnvVars recursively expands environment variables in string fields of a PollerConfig
func expandPollerEnvVars(pollerCfg PollerConfig) PollerConfig {
	// Using reflection to iterate over struct fields
	val := reflect.ValueOf(&pollerCfg).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			originalValue := field.String()
			expandedValue := os.ExpandEnv(originalValue)
			if field.CanSet() {
				field.SetString(expandedValue)
			}
		}
	}
	return pollerCfg
}
