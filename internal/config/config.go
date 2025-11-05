package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config holds the application configuration
type Config struct {
	Watchdog WatchdogConfig `mapstructure:"watchdog"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// WatchdogConfig holds the watchdog-specific configuration
type WatchdogConfig struct {
	Namespaces       []string          `mapstructure:"namespaces"`
	LabelSelectors   map[string]string `mapstructure:"labelSelectors"`
	ScheduleInterval time.Duration     `mapstructure:"scheduleInterval"`
	MaxPodLifetime   time.Duration     `mapstructure:"maxPodLifetime"`
	DryRun           bool              `mapstructure:"dryRun"`
}

// LoggingConfig holds the logging configuration
type LoggingConfig struct {
	Mode  string `mapstructure:"mode"`
	Level string `mapstructure:"level"`
}

// NewConfig loads the configuration from the config file
func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("configs")

	// Set default values
	viper.SetDefault("watchdog.scheduleInterval", "10m")
	viper.SetDefault("watchdog.maxPodLifetime", "24h")
	viper.SetDefault("watchdog.dryRun", false)
	viper.SetDefault("logging.mode", "production")
	viper.SetDefault("logging.level", "info")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Module provides the configuration as a dependency
var Module = fx.Options(
	fx.Provide(NewConfig),
)

// ProvideConfigForTest provides a config for testing
func ProvideConfigForTest() *Config {
	return &Config{
		Watchdog: WatchdogConfig{
			Namespaces: []string{"default"},
			LabelSelectors: map[string]string{
				"app": "test",
			},
			ScheduleInterval: 10 * time.Minute,
			MaxPodLifetime:   24 * time.Hour,
			DryRun:           true,
		},
		Logging: LoggingConfig{
			Mode:  "development",
			Level: "debug",
		},
	}
}
