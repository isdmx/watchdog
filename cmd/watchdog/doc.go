// Package main contains the entry point for the watchdog application.
//
// The watchdog is a Kubernetes pod lifetime management tool that monitors
// pods in specified namespaces and terminates those that exceed a configured
// maximum lifetime. It provides a simple command-line interface to start
// the monitoring service.
//
// The application uses the Uber fx framework for dependency injection and
// follows a modular architecture with separate packages for configuration,
// logging, Kubernetes client interaction, monitoring, and HTTP server functionality.
//
// Features:
//   - Configurable maximum pod lifetime enforcement
//   - Namespace and label selector filtering
//   - Dry-run mode for testing
//   - Prometheus metrics collection
//   - Health and readiness endpoints
package main
