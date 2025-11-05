package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/config"
	"github.com/isdmx/watchdog/handlers"
	"github.com/isdmx/watchdog/monitoring"
)

func main() {
	app := fx.New(
		// Configuration module
		config.Module,

		// Logging module
		config.LoggingModule,

		// Kubernetes client module
		config.KubernetesModule,

		// Health check handler
		handlers.HealthCheckModule,

		// HTTP server
		fx.Provide(newHTTPServer),

		// Monitoring module
		monitoring.MonitoringModule,

		// Start the application
		fx.Invoke(startMonitoring),
	)

	// Start the application
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shut down the application gracefully
	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}

// newHTTPServer creates a new HTTP server with health check endpoints
func newHTTPServer(healthHandler *handlers.HealthHandler) *http.Server {
	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

// startMonitoring starts the monitoring process
func startMonitoring(lc fx.Lifecycle, server *http.Server, pm *monitoring.PodMonitor, cfg *config.Config, logger *zap.SugaredLogger) {
	// Add the HTTP server to the lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server on :8080")
			// Start the HTTP server in a goroutine
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		},
	})

	// Start the monitoring in the background
	go monitoring.StartMonitoring(pm, cfg, logger)

	// Add monitoring stop to the lifecycle
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down monitoring process")
			pm.Stop()
			return nil
		},
	})
}