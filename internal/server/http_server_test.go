package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"
)

func TestNewHttpServer(t *testing.T) {
	t.Run("creates HTTP server successfully", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()
		sugaredLogger := logger.Sugar()

		cfg := &config.Config{
			Http: config.HttpConfig{
				Addr:         ":8080",
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			},
		}

		lc := fxtest.NewLifecycle(t)
		server := NewHTTPServer(lc, sugaredLogger, cfg)

		require.NotNil(t, server)
		require.NotNil(t, server.server)
		require.Equal(t, cfg.Http.Addr, server.server.Addr)
		require.Equal(t, cfg.Http.ReadTimeout, server.server.ReadTimeout)
		require.Equal(t, cfg.Http.WriteTimeout, server.server.WriteTimeout)
	})
}

func TestRegisterRoutes(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugaredLogger := logger.Sugar()

	cfg := &config.Config{
		Http: config.HttpConfig{
			Addr:         ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	lc := fxtest.NewLifecycle(t)
	server := NewHTTPServer(lc, sugaredLogger, cfg)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	// Create a test server with the registered routes
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	t.Run("healthz endpoint", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", testServer.URL+"/healthz", http.NoBody)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "OK", strings.TrimSpace(string(body)))
	})

	t.Run("readyz endpoint", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", testServer.URL+"/readyz", http.NoBody)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "OK", strings.TrimSpace(string(body)))
	})

	t.Run("metrics endpoint", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", testServer.URL+"/metrics", http.NoBody)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		// The metrics endpoint should return prometheus metrics
		require.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
	})
}

func TestHttpServerMethods(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugaredLogger := logger.Sugar()

	cfg := &config.Config{
		Http: config.HttpConfig{
			Addr:         ":0", // Use random port
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	lc := fxtest.NewLifecycle(t)
	server := NewHTTPServer(lc, sugaredLogger, cfg)

	t.Run("start and shutdown server", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Start the server
		err := server.Start(ctx)
		require.NoError(t, err)

		// Give the server a moment to start
		time.Sleep(100 * time.Millisecond)

		// The exact address depends on the random port,
		// so we'll just check that shutdown works without error
		err = server.Shutdown(context.Background())
		require.NoError(t, err)
	})
}

func TestHttpHandlerFunctions(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugaredLogger := logger.Sugar()

	server := &HTTPServer{
		logger: sugaredLogger,
	}

	t.Run("healthz handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/healthz", http.NoBody)
		w := httptest.NewRecorder()

		server.healthz(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "text/plain", w.Header().Get("Content-Type"))
		require.Equal(t, "OK", w.Body.String())
	})

	t.Run("readyz handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/readyz", http.NoBody)
		w := httptest.NewRecorder()

		server.readyz(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "text/plain", w.Header().Get("Content-Type"))
		require.Equal(t, "OK", w.Body.String())
	})

	t.Run("metrics handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/metrics", http.NoBody)
		w := httptest.NewRecorder()

		server.metrics(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		// The content type should be text/plain for prometheus metrics
		require.Contains(t, w.Header().Get("Content-Type"), "text/plain")
	})
}
