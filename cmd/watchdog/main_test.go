package main

import (
	"context"
	"testing"
	"time"

	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/isdmx/watchdog/internal/config"
	"github.com/isdmx/watchdog/internal/handlers"
	"github.com/isdmx/watchdog/internal/monitoring"
)

// Create a fake Kubernetes client for testing
func NewFakeKubernetesClient() kubernetes.Interface {
	clientset := fake.NewSimpleClientset()
	return clientset
}

// TestNewHTTPServer tests the newHTTPServer function
func TestNewHTTPServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	healthHandler := handlers.NewHealthHandler(logger.Sugar())

	server := newHTTPServer(healthHandler)

	if server.Addr != ":8080" {
		t.Errorf("Expected server address to be :8080, got %s", server.Addr)
	}

	if server.ReadHeaderTimeout != 20*time.Second {
		t.Errorf("Expected ReadHeaderTimeout to be 20s, got %v", server.ReadHeaderTimeout)
	}

	// Since we don't have the full mux setup in isolation, we can at least verify the server was created
	if server.Handler == nil {
		t.Error("Expected server handler to be set")
	}
}

// TestStartMonitoringFunction tests the startMonitoring function by using Fx testing utilities
func TestStartMonitoringFunction(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	// Create a test context
	lc := fxtest.NewLifecycle(t)

	// Create required dependencies
	testConfig := &config.Config{
		Watchdog: config.WatchdogConfig{
			Namespaces:       []string{"default"},
			LabelSelectors:   map[string]string{"app": "test"},
			MaxPodLifetime:   time.Hour,
			DryRun:           true,
			ScheduleInterval: time.Minute,
		},
	}

	// Create a mock pod monitor or a real one with a fake client
	clientset := NewFakeKubernetesClient()
	podMonitor := monitoring.NewPodMonitor(clientset, testConfig, logger.Sugar())

	// Create a dummy HTTP server
	healthHandler := handlers.NewHealthHandler(logger.Sugar())
	server := newHTTPServer(healthHandler)

	// Test that startMonitoring completes without error
	// Since startMonitoring runs server and monitoring in goroutines,
	// we just need to make sure it doesn't panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("startMonitoring panicked: %v", r)
			}
		}()

		// This will run the server in a goroutine, so we should stop it afterwards
		startMonitoring(lc, server, podMonitor, testConfig, logger.Sugar())
	}()

	// Stop the lifecycle to clean up
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := lc.Stop(ctx); err != nil {
		t.Errorf("Failed to stop lifecycle: %v", err)
	}
}
