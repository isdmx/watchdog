package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoading(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	
	configContent := `
watchdog:
  namespaces:
    - "test-namespace"
    - "another-namespace"
  labelSelectors:
    app: "test-app"
    version: "v1"
  scheduleInterval: "5m"
  maxPodLifetime: "2h"
  dryRun: true

logging:
  mode: "development"
  level: "debug"
`
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Change to temp dir to read the config file
	oldDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Load the config
	cfg, err := NewConfig()
	require.NoError(t, err)

	// Verify the config values
	require.NotNil(t, cfg)
	
	assert.Equal(t, []string{"test-namespace", "another-namespace"}, cfg.Watchdog.Namespaces)
	assert.Equal(t, "test-app", cfg.Watchdog.LabelSelectors["app"])
	assert.Equal(t, "v1", cfg.Watchdog.LabelSelectors["version"])
	assert.Equal(t, 5*time.Minute, cfg.Watchdog.ScheduleInterval)
	assert.Equal(t, 2*time.Hour, cfg.Watchdog.MaxPodLifetime)
	assert.True(t, cfg.Watchdog.DryRun)
	assert.Equal(t, "development", cfg.Logging.Mode)
	assert.Equal(t, "debug", cfg.Logging.Level)
}

func TestConfigWithDefaults(t *testing.T) {
	// Test with default values when config is not present
	cfg := ProvideConfigForTest()
	
	require.NotNil(t, cfg)
	assert.Equal(t, []string{"default"}, cfg.Watchdog.Namespaces)
	assert.Equal(t, "test", cfg.Watchdog.LabelSelectors["app"])
	assert.Equal(t, 10*time.Minute, cfg.Watchdog.ScheduleInterval)
	assert.Equal(t, 24*time.Hour, cfg.Watchdog.MaxPodLifetime)
	assert.True(t, cfg.Watchdog.DryRun)
	assert.Equal(t, "development", cfg.Logging.Mode)
	assert.Equal(t, "debug", cfg.Logging.Level)
}