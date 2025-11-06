package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/isdmx/watchdog/internal/config"
	"github.com/isdmx/watchdog/internal/monitoring"
)

// MockPodMonitor is a mock implementation of monitoring.PodMonitor
type MockPodMonitor struct {
	mock.Mock
}

func (m *MockPodMonitor) MonitorAndCleanup() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewWatchdogServer(t *testing.T) {
	t.Run("creates watchdog server successfully", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()
		configObj := &config.Config{
			Watchdog: config.WatchdogConfig{
				ScheduleInterval: 10 * time.Minute,
				MaxPodLifetime:   1 * time.Hour,
			},
		}
		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		podMonitor := monitoring.NewPodMonitor(clientset, configObj, sugaredLogger)

		lc := fxtest.NewLifecycle(t)
		wdServer := NewWatchdogServer(lc, podMonitor, sugaredLogger, configObj)

		require.NotNil(t, wdServer)
		require.Equal(t, podMonitor, wdServer.pm)
		require.Equal(t, configObj, wdServer.config)
		require.NotNil(t, wdServer.stopChannel)
	})
}

func TestWatchdogServerStart(t *testing.T) {
	t.Run("starts monitoring successfully", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()
		configObj := &config.Config{
			Watchdog: config.WatchdogConfig{
				ScheduleInterval: 100 * time.Millisecond, // Fast interval for testing
				MaxPodLifetime:   1 * time.Hour,
			},
		}
		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		podMonitor := monitoring.NewPodMonitor(clientset, configObj, sugaredLogger)

		lc := fxtest.NewLifecycle(t)
		wdServer := NewWatchdogServer(lc, podMonitor, sugaredLogger, configObj)

		ctx := context.Background()
		err := wdServer.Start(ctx)
		require.NoError(t, err)

		// Give a little time for the goroutine to start
		time.Sleep(150 * time.Millisecond)

		// Shutdown to clean up
		err = wdServer.Shutdown(ctx)
		require.NoError(t, err)
	})
}

func TestWatchdogServerShutdown(t *testing.T) {
	t.Run("shuts down monitoring properly", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()
		configObj := &config.Config{
			Watchdog: config.WatchdogConfig{
				ScheduleInterval: 10 * time.Minute, // Slow interval so it doesn't run during test
				MaxPodLifetime:   1 * time.Hour,
			},
		}
		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		podMonitor := monitoring.NewPodMonitor(clientset, configObj, sugaredLogger)

		lc := fxtest.NewLifecycle(t)
		wdServer := NewWatchdogServer(lc, podMonitor, sugaredLogger, configObj)

		ctx := context.Background()
		err := wdServer.Start(ctx)
		require.NoError(t, err)

		// Shutdown the server
		err = wdServer.Shutdown(ctx)
		require.NoError(t, err)

		// Verify the stop channel was closed
		select {
		case _, ok := <-wdServer.stopChannel:
			// If the channel is closed, ok will be false
			require.False(t, ok)
		default:
			// Channel should be closed, so this shouldn't happen
			require.Fail(t, "stopChannel was not closed")
		}
	})
}
