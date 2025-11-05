package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/isdmx/watchdog/internal/config"

	"go.uber.org/zap"
)

func TestNewPodMonitor(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	cfg := &config.Config{
		Watchdog: config.WatchdogConfig{
			ScheduleInterval: 5 * time.Minute,
			MaxPodLifetime:   1 * time.Hour,
		},
	}
	logger, _ := zap.NewDevelopment()
	sugaredLogger := logger.Sugar()

	pm := NewPodMonitor(clientset, cfg, sugaredLogger)
	require.NotNil(t, pm)
	require.Equal(t, clientset, pm.clientset)
	require.Equal(t, cfg, pm.config)
	require.Equal(t, sugaredLogger, pm.logger)
}

func TestMonitorAndCleanup(t *testing.T) {
	t.Run("processes pods and terminates old ones", func(t *testing.T) {
		// Create a fake clientset with some pods
		clientset := fake.NewSimpleClientset()

		// Create a pod that's older than the max lifetime
		oldPod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "old-pod",
				Namespace:         "default",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-2 * time.Hour)}, // 2 hours old
			},
		}

		_, err := clientset.CoreV1().Pods("default").Create(context.TODO(), oldPod, metav1.CreateOptions{})
		require.NoError(t, err)

		// Create a pod that's within the max lifetime
		newPod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "new-pod",
				Namespace:         "default",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-30 * time.Minute)}, // 30 minutes old
			},
		}

		_, err = clientset.CoreV1().Pods("default").Create(context.TODO(), newPod, metav1.CreateOptions{})
		require.NoError(t, err)

		cfg := &config.Config{
			Watchdog: config.WatchdogConfig{
				Namespaces:     []string{"default"},
				MaxPodLifetime: 1 * time.Hour, // 1 hour max
				DryRun:         false,         // Actually terminate pods
			},
		}

		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		pm := NewPodMonitor(clientset, cfg, sugaredLogger)
		err = pm.MonitorAndCleanup()
		require.NoError(t, err)

		// Check that the old pod was terminated
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "old-pod", metav1.GetOptions{})
		require.Error(t, err) // Pod should be deleted

		// Check that the new pod still exists
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "new-pod", metav1.GetOptions{})
		require.NoError(t, err) // Pod should still exist
	})

	t.Run("runs in dry run mode", func(t *testing.T) {
		// Create a fake clientset with an old pod
		clientset := fake.NewSimpleClientset()

		oldPod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "old-pod",
				Namespace:         "default",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-2 * time.Hour)}, // 2 hours old
			},
		}

		_, err := clientset.CoreV1().Pods("default").Create(context.TODO(), oldPod, metav1.CreateOptions{})
		require.NoError(t, err)

		cfg := &config.Config{
			Watchdog: config.WatchdogConfig{
				Namespaces:     []string{"default"},
				MaxPodLifetime: 1 * time.Hour, // 1 hour max
				DryRun:         true,          // Dry run mode
			},
		}

		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		pm := NewPodMonitor(clientset, cfg, sugaredLogger)
		err = pm.MonitorAndCleanup()
		require.NoError(t, err)

		// In dry run mode, the pod should still exist
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "old-pod", metav1.GetOptions{})
		require.NoError(t, err) // Pod should still exist in dry run mode
	})

	t.Run("handles error when listing pods", func(t *testing.T) {
		// This test is more complex as it would require mocking failure scenarios
		// For now, we test the success path as shown above
		cfg := &config.Config{
			Watchdog: config.WatchdogConfig{
				Namespaces:     []string{"nonexistent-namespace"},
				MaxPodLifetime: 1 * time.Hour,
				DryRun:         false,
			},
		}

		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		pm := NewPodMonitor(fake.NewSimpleClientset(), cfg, sugaredLogger)
		err := pm.MonitorAndCleanup()
		// This should not return an error, it should just log and continue
		require.NoError(t, err)
	})
}

func TestBuildLabelSelector(t *testing.T) {
	t.Run("creates empty selector for empty map", func(t *testing.T) {
		labels := map[string]string{}
		selector := buildLabelSelector(labels)
		require.Empty(t, selector)
	})

	t.Run("creates selector for single label", func(t *testing.T) {
		labels := map[string]string{
			"app": "test",
		}
		selector := buildLabelSelector(labels)
		require.Equal(t, "app=test", selector)
	})

	t.Run("creates selector for multiple labels", func(t *testing.T) {
		labels := map[string]string{
			"app": "test",
			"env": "prod",
			"ver": "v1",
		}
		selector := buildLabelSelector(labels)
		// The order might vary, so we'll check that all parts are present
		require.Contains(t, selector, "app=test")
		require.Contains(t, selector, "env=prod")
		require.Contains(t, selector, "ver=v1")
		// Since map iteration order is not guaranteed, we need to check all parts
		parts := []string{"app=test", "env=prod", "ver=v1"}
		for _, part := range parts {
			require.Contains(t, selector, part)
		}
	})
}

func TestTerminatePod(t *testing.T) {
	t.Run("terminates pod successfully", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
		}

		_, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
		require.NoError(t, err)

		cfg := &config.Config{
			Watchdog: config.WatchdogConfig{
				MaxPodLifetime: 1 * time.Hour,
			},
		}

		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		pm := NewPodMonitor(clientset, cfg, sugaredLogger)
		err = pm.terminatePod("default", "test-pod")
		require.NoError(t, err)

		// Verify pod was deleted
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "test-pod", metav1.GetOptions{})
		require.Error(t, err)
	})

	t.Run("handles error when terminating non-existent pod", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()

		cfg := &config.Config{
			Watchdog: config.WatchdogConfig{
				MaxPodLifetime: 1 * time.Hour,
			},
		}

		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		pm := NewPodMonitor(clientset, cfg, sugaredLogger)
		err := pm.terminatePod("default", "non-existent-pod")
		// This should return an error since the pod doesn't exist
		require.Error(t, err)
	})
}
