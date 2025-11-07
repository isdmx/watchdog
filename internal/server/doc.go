// Package server provides HTTP server functionality and the watchdog server implementation.
//
// This package contains two main server implementations:
//   - HTTPServer: A standard HTTP server that provides health check endpoints
//     (/healthz, /readyz) and Prometheus metrics (/metrics)
//   - WatchdogServer: A specialized server that runs the pod monitoring logic
//     on a configurable schedule
//
// The HTTP server supports:
//   - Configurable address and timeouts
//   - Health and readiness check endpoints for container orchestration platforms
//   - Prometheus metrics endpoint for monitoring and alerting
//   - Lifecycle management via the fx framework
//
// The Watchdog server provides:
//   - Periodic pod monitoring based on configuration schedule
//   - Integration with the monitoring package for pod termination
//   - Graceful startup and shutdown handling
//   - Lifecycle management to ensure proper cleanup
//
// Both servers implement the Server interface which defines common Start and
// Shutdown methods, enabling consistent lifecycle management across different
// server types. The package is designed to work with the Uber fx framework's
// lifecycle hooks for proper initialization and cleanup.
package server
