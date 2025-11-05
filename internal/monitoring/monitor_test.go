package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/isdmx/watchdog/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPodMonitor(t *testing.T) {
	// Create a fake kubernetes client for testing
	fakeClient := fake.NewSimpleClientset()

	// Create a test config
	testConfig := &config.Config{
		Watchdog: config.WatchdogConfig{
			Namespaces: []string{"test-namespace"},
			LabelSelectors: map[string]string{
				"app": "test",
			},
			ScheduleInterval: 10 * time.Minute,
			MaxPodLifetime:   1 * time.Hour,
			DryRun:           true, // Set to true for safety in tests
		},
	}

	// Create logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() {
		// Best effort to sync logger
		_ = logger.Sync()
	}()
	sugarLogger := logger.Sugar()

	// Create pod monitor - the fake client implements the kubernetes.Interface
	podMonitor := NewPodMonitor(fakeClient, testConfig, sugarLogger)

	t.Run("MonitorAndCleanup with no pods", func(t *testing.T) {
		err := podMonitor.MonitorAndCleanup()
		assert.NoError(t, err)
	})

	t.Run("MonitorAndCleanup with old pod", func(t *testing.T) {
		// Create a pod that is older than the max lifetime
		oldPod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "old-pod",
				Namespace:         "test-namespace",
				Labels:            map[string]string{"app": "test"},
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-2 * time.Hour)}, // 2 hours old
			},
		}

		_, err := fakeClient.CoreV1().Pods("test-namespace").Create(context.TODO(), oldPod, metav1.CreateOptions{})
		require.NoError(t, err)

		err = podMonitor.MonitorAndCleanup()
		require.NoError(t, err)
	})

	t.Run("BuildLabelSelector", func(t *testing.T) {
		labels := map[string]string{
			"app":     "test",
			"version": "v1",
		}

		selector := buildLabelSelector(labels)

		// Since map iteration order is not guaranteed, we'll check that both key-value pairs are present
		assert.Contains(t, selector, "app=test")
		assert.Contains(t, selector, "version=v1")
		assert.Contains(t, selector, ",")
	})
}
