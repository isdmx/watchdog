package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"
)

var _ Server = (*HttpServer)(nil)

// HttpServer manages health check endpoints
type HttpServer struct {
	server *http.Server
	logger *zap.SugaredLogger
}

// NewHttpServer creates a new health handler
func NewHttpServer(lc fx.Lifecycle, logger *zap.SugaredLogger, cfg *config.Config) *HttpServer {
	mux := http.NewServeMux()

	server := &HttpServer{
		logger: logger,
		server: &http.Server{
			Handler:      mux,
			Addr:         cfg.Http.Addr,
			ReadTimeout:  cfg.Http.ReadTimeout,
			WriteTimeout: cfg.Http.WriteTimeout,
		},
	}

	server.RegisterRoutes(mux)
	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Shutdown,
	})
	return server
}

// RegisterRoutes registers health check routes
func (h *HttpServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/metrics", h.metrics)
}

// healthz endpoint - checks if the service is healthy
func (h *HttpServer) healthz(w http.ResponseWriter, _ *http.Request) {
	// In a real implementation, you might check connectivity to external services
	// For now, we'll just return OK
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.Error("Failed to write healthz response", "error", err)
	}
}

// readyz endpoint - checks if the service is ready to serve requests
func (h *HttpServer) readyz(w http.ResponseWriter, _ *http.Request) {
	// In a real implementation, you might check if other dependencies are ready
	// For now, we'll just return OK
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.Error("Failed to write readyz response", "error", err)
	}
}

// metrics endpoint - returns Prometheus metrics
func (*HttpServer) metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func (h *HttpServer) Start(ctx context.Context) error {
	h.logger.Info("Starting HTTP server")

	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Done()
			h.logger.Error("Failed to start HTTP server", "error", err)
		}
	}(h.server)

	return nil
}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	h.logger.Info("Shutting down HTTP server")
	return h.server.Shutdown(ctx)
}
