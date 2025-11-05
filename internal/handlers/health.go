package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// HealthHandler manages health check endpoints
type HealthHandler struct {
	logger *zap.SugaredLogger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger *zap.SugaredLogger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/metrics", h.metrics)
}

// healthz endpoint - checks if the service is healthy
func (h *HealthHandler) healthz(w http.ResponseWriter, _ *http.Request) {
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
func (h *HealthHandler) readyz(w http.ResponseWriter, _ *http.Request) {
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
func (*HealthHandler) metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// HealthCheckModule provides the health check handler as a dependency
var HealthCheckModule = fx.Options(
	fx.Provide(NewHealthHandler),
	fx.Invoke(RegisterHealthRoutes),
)

// RegisterHealthRoutes registers the health check routes
func RegisterHealthRoutes(healthHandler *HealthHandler, mux *http.ServeMux) {
	healthHandler.RegisterRoutes(mux)
}
