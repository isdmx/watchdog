package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"
)

var _ Server = (*HTTPServer)(nil)

// HTTPServer manages health check endpoints
type HTTPServer struct {
	server *http.Server
	logger *zap.SugaredLogger
}

// NewHTTPServer creates a new health handler
func NewHTTPServer(lc fx.Lifecycle, logger *zap.SugaredLogger, cfg *config.Config) *HTTPServer {
	mux := http.NewServeMux()

	server := &HTTPServer{
		logger: logger.Named("HTTPServer"),
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
func (h *HTTPServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/metrics", h.metrics)
}

// healthz endpoint - checks if the service is healthy
func (h *HTTPServer) healthz(w http.ResponseWriter, _ *http.Request) {
	// In a real implementation, you might check connectivity to external services
	// For now, we'll just return OK
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.Errorw("Failed to write healthz response", "error", err)
	}
}

// readyz endpoint - checks if the service is ready to serve requests
func (h *HTTPServer) readyz(w http.ResponseWriter, _ *http.Request) {
	// In a real implementation, you might check if other dependencies are ready
	// For now, we'll just return OK
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.Errorw("Failed to write readyz response", "error", err)
	}
}

// metrics endpoint - returns Prometheus metrics
func (*HTTPServer) metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func (h *HTTPServer) Start(ctx context.Context) error {
	h.logger.Info("Starting HTTP server")

	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Done()
			h.logger.Errorw("Failed to start HTTP server", "error", err)
		}
	}(h.server)

	return nil
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	h.logger.Info("Shutting down HTTP server")
	return h.server.Shutdown(ctx)
}
