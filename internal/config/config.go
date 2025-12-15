package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Watchdog WatchdogConfig `mapstructure:"watchdog"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	HTTP     HTTPConfig     `mapstructure:"http"`
}

// HTTPConfig holds the healthcheck-specific configuration
type HTTPConfig struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
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

const (
	defaultHTTPAddr         = ":8080"
	defaultReadTimeout      = 5 * time.Second
	defaultWriteTimeout     = 10 * time.Second
	defaultScheduleInterval = 10 * time.Minute
	defaultMaxPodLifetime   = 1 * time.Hour
	defaultDryRun           = false
	defaultLogLevel         = "info"
	defaultLogMode          = "production"
)

// NewConfig loads the configuration from the config file
func NewConfig() (*Config, error) {
	viper.SetOptions(viper.KeyDelimiter("::")) // because labelSelectors may contain `.`
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")

	// Set default values
	viper.SetDefault("http::addr", defaultHTTPAddr)
	viper.SetDefault("http::readTimeout", defaultReadTimeout)
	viper.SetDefault("http::writeTimeout", defaultWriteTimeout)
	viper.SetDefault("watchdog::scheduleInterval", defaultScheduleInterval)
	viper.SetDefault("watchdog::maxPodLifetime", defaultMaxPodLifetime)
	viper.SetDefault("watchdog::dryRun", defaultDryRun)
	viper.SetDefault("logging::mode", defaultLogMode)
	viper.SetDefault("logging::level", defaultLogLevel)

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
