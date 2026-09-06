package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// Config holds the root configuration for adxctl
type Config struct {
	Nats    NatsConfig              `mapstructure:"nats" yaml:"nats"`
	Server  NatsConfig              `mapstructure:"server" yaml:"server"` // Backwards compatibility fallback
	Pollers map[string]PollerConfig `mapstructure:"pollers" yaml:"pollers"`
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

	// Expand environment variables in poller configurations
	for name, pollerCfg := range cfg.Pollers {
		expandedPollerCfg := expandPollerEnvVars(pollerCfg)
		cfg.Pollers[name] = expandedPollerCfg
	}

	// If nats was not explicitly defined in config file but server was, copy from server
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
