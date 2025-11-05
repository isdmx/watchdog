package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHealthHandler(t *testing.T) {
	// Create logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() {
		// Best effort to sync logger
		_ = logger.Sync()
	}()
	sugarLogger := logger.Sugar()

	// Create health handler
	healthHandler := NewHealthHandler(sugarLogger)

	// Create a mux and register routes
	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)

	t.Run("Healthz endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/healthz", http.NoBody)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("Readyz endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/readyz", http.NoBody)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("Metrics endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/metrics", http.NoBody)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// The metrics endpoint should return Prometheus metrics format
		assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")
	})
}
