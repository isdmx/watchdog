package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("loads config successfully", func(t *testing.T) {
		// Create a temporary config file
		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(cfgPath, []byte(`
watchdog:
  namespaces: ["default", "kube-system"]
  labelSelectors:
    app: test
  scheduleInterval: 5m
  maxPodLifetime: 30m
  dryRun: true
http:
  addr: ":9090"
  readTimeout: 10s
  writeTimeout: 20s
logging:
  mode: development
  level: debug
`), 0600)
		require.NoError(t, err)

		// Change to the temp directory to find the config
		origDir, err := os.Getwd()
		require.NoError(t, err)
		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = os.Chdir(origDir)
		})

		config, err := NewConfig()
		require.NoError(t, err)
		require.NotNil(t, config)

		// Check watchdog config
		require.Equal(t, []string{"default", "kube-system"}, config.Watchdog.Namespaces)
		require.Equal(t, map[string]string{"app": "test"}, config.Watchdog.LabelSelectors)
		require.Equal(t, 5*time.Minute, config.Watchdog.ScheduleInterval)
		require.Equal(t, 30*time.Minute, config.Watchdog.MaxPodLifetime)
		require.True(t, config.Watchdog.DryRun)

		// Check http config
		require.Equal(t, ":9090", config.Http.Addr)
		require.Equal(t, 10*time.Second, config.Http.ReadTimeout)
		require.Equal(t, 20*time.Second, config.Http.WriteTimeout)

		// Check logging config
		require.Equal(t, "development", config.Logging.Mode)
		require.Equal(t, "debug", config.Logging.Level)
	})

	t.Run("uses default values when config is missing", func(t *testing.T) {
		// Create a temp directory without config file
		tmpDir := t.TempDir()

		origDir, err := os.Getwd()
		require.NoError(t, err)
		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = os.Chdir(origDir)
		})

		config, err := NewConfig()
		// This should fail because it can't find the config file
		require.Error(t, err)
		require.Nil(t, config)
	})

	t.Run("handles invalid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(cfgPath, []byte("invalid yaml"), 0600)
		require.NoError(t, err)

		origDir, err := os.Getwd()
		require.NoError(t, err)
		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = os.Chdir(origDir)
		})

		config, err := NewConfig()
		require.Error(t, err)
		require.Nil(t, config)
	})
}

func TestConfigDefaults(t *testing.T) {
	// Test with default values when config file exists but doesn't specify all values
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(cfgPath, []byte(`watchdog: {}`), 0600)
	require.NoError(t, err)

	origDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	config, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Check defaults
	require.Equal(t, defaultHttpAddr, config.Http.Addr)
	require.Equal(t, defaultReadTimeout, config.Http.ReadTimeout)
	require.Equal(t, defaultWriteTimeout, config.Http.WriteTimeout)
	require.Equal(t, defaultScheduleInterval, config.Watchdog.ScheduleInterval)
	require.Equal(t, defaultMaxPodLifetime, config.Watchdog.MaxPodLifetime)
	require.Equal(t, defaultDryRun, config.Watchdog.DryRun)
	require.Equal(t, defaultLogMode, config.Logging.Mode)
	require.Equal(t, defaultLogLevel, config.Logging.Level)
}
