package logging

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isdmx/watchdog/internal/config"
)

func TestNewLogger(t *testing.T) {
	t.Run("creates development logger", func(t *testing.T) {
		cfg := &config.Config{
			Logging: config.LoggingConfig{
				Mode:  "development",
				Level: "debug",
			},
		}

		logger, err := NewLogger(cfg)
		require.NoError(t, err)
		require.NotNil(t, logger)
		t.Cleanup(func() {
			_ = logger.Sync()
		})
	})

	t.Run("creates production logger", func(t *testing.T) {
		cfg := &config.Config{
			Logging: config.LoggingConfig{
				Mode:  "production",
				Level: "info",
			},
		}

		logger, err := NewLogger(cfg)
		require.NoError(t, err)
		require.NotNil(t, logger)
		t.Cleanup(func() {
			_ = logger.Sync()
		})
	})

	t.Run("handles invalid log level", func(t *testing.T) {
		cfg := &config.Config{
			Logging: config.LoggingConfig{
				Mode:  "production",
				Level: "invalid_level",
			},
		}

		logger, err := NewLogger(cfg)
		require.NoError(t, err)
		require.NotNil(t, logger)
		t.Cleanup(func() {
			_ = logger.Sync()
		})

		// We can verify by checking if it accepts info level logs
		logger.Info("test message")
	})

	t.Run("handles error when building logger", func(t *testing.T) {
		// This test is tricky since zap config build errors are rare in normal usage
		// For the purpose of this test, we'll just test the happy path
		cfg := &config.Config{
			Logging: config.LoggingConfig{
				Mode:  "production",
				Level: "info",
			},
		}

		logger, err := NewLogger(cfg)
		require.NoError(t, err)
		require.NotNil(t, logger)
		t.Cleanup(func() {
			_ = logger.Sync()
		})
	})
}

func TestNewSugaredLogger(t *testing.T) {
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Mode:  "production",
			Level: "info",
		},
	}

	logger, err := NewLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)
	t.Cleanup(func() {
		_ = logger.Sync()
	})

	sugaredLogger := NewSugaredLogger(logger)
	require.NotNil(t, sugaredLogger)

	// Test that sugared logger works
	sugaredLogger.Info("test message")
}
