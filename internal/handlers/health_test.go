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

func TestRegisterHealthRoutes(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	healthHandler := NewHealthHandler(logger.Sugar())
	mux := http.NewServeMux()

	// Test that the function doesn't panic and registers routes
	RegisterHealthRoutes(healthHandler, mux)

	// Verify the routes were registered by checking that requests don't return 404
	testCases := []struct {
		path string
	}{
		{"/healthz"},
		{"/readyz"},
		{"/metrics"},
	}

	for _, tc := range testCases {
		t.Run("route_"+tc.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, http.NoBody)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			// Check that it's not a 404 (route was registered)
			if w.Code == http.StatusNotFound {
				t.Errorf("Route %s was not registered, got 404", tc.path)
			}
		})
	}
}
