package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestTerminatePod(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	cfg := &config.Config{
		Watchdog: config.WatchdogConfig{
			Namespaces:       []string{"test-namespace"},
			LabelSelectors:   map[string]string{"app": "test"},
			MaxPodLifetime:   time.Hour,
			DryRun:           false,
			ScheduleInterval: time.Minute,
		},
	}

	podMonitor := NewPodMonitor(clientset, cfg, logger.Sugar())

	// Create a test pod
	testPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
	}

	_, err := clientset.CoreV1().Pods("test-namespace").Create(context.TODO(), testPod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test pod: %v", err)
	}

	// Test terminating the pod
	err = podMonitor.terminatePod("test-namespace", "test-pod")
	if err != nil {
		t.Errorf("terminatePod failed: %v", err)
	}

	// Verify the pod is deleted
	_, err = clientset.CoreV1().Pods("test-namespace").Get(context.TODO(), "test-pod", metav1.GetOptions{})
	if err == nil {
		t.Error("Expected pod to be deleted, but it still exists")
	}
}

func TestStartMonitoring(_ *testing.T) {
	clientset := fake.NewSimpleClientset()
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	cfg := &config.Config{
		Watchdog: config.WatchdogConfig{
			Namespaces:       []string{"test-namespace"},
			LabelSelectors:   map[string]string{"app": "test"},
			MaxPodLifetime:   500 * time.Millisecond, // Short lifetime for testing
			DryRun:           true,                   // Use dry run to avoid actual deletions
			ScheduleInterval: 100 * time.Millisecond, // Short interval for testing
		},
	}

	podMonitor := NewPodMonitor(clientset, cfg, logger.Sugar())

	// Start monitoring in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		StartMonitoring(podMonitor, cfg, logger.Sugar())
	}()

	// Let it run for a short time
	time.Sleep(250 * time.Millisecond)

	// Stop the monitoring
	podMonitor.Stop()

	wg.Wait() // Wait for StartMonitoring to return
}

func TestStop(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	cfg := &config.Config{
		Watchdog: config.WatchdogConfig{
			Namespaces:       []string{"test-namespace"},
			LabelSelectors:   map[string]string{"app": "test"},
			MaxPodLifetime:   time.Hour,
			DryRun:           true,
			ScheduleInterval: time.Minute,
		},
	}

	podMonitor := NewPodMonitor(clientset, cfg, logger.Sugar())

	// Initially, the stop channel should be open
	select {
	case _, ok := <-podMonitor.stopChannel:
		if ok {
			// This shouldn't happen - if we got a value with ok=true,
			// that means there was a value in the closed channel, which is unexpected
		} else {
			// Channel was closed and we received the zero value, which is unexpected initially
			t.Error("Expected stopChannel to be open initially")
		}
	default:
		// Channel is open but no value available for reading, which is expected initially
	}

	// Call Stop
	podMonitor.Stop()

	// After Stop, the channel should be closed
	select {
	case _, ok := <-podMonitor.stopChannel:
		if ok {
			t.Error("Expected stopChannel to be closed after Stop() call")
		}
	default:
		t.Error("Expected stopChannel to be closed and readable after Stop() call")
	}
}

func TestMonitoringModule(t *testing.T) {
	// Test that the module is defined
	if MonitoringModule == nil {
		t.Error("MonitoringModule should not be nil")
	}
}
