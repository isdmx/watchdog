// Package monitoring provides pod monitoring and cleanup functionality.
//
// This package implements the core business logic for the watchdog application,
// specifically monitoring pods in Kubernetes and terminating those that exceed
// the configured maximum lifetime. It provides the PodMonitor type which handles
// the monitoring operations and integrates with various Kubernetes API resources.
//
// Key features of this package:
//   - Pod lifetime monitoring based on creation timestamp
//   - Namespace and label selector filtering for targeted monitoring
//   - Dry-run mode for safe testing of monitoring policies
//   - Prometheus metrics collection for monitoring operations
//   - Configurable monitoring intervals and maximum pod lifetimes
//
// The package includes:
//   - PodMonitor: Main monitoring type that handles pod inspection and termination
//   - Metrics: Prometheus metrics for tracking monitoring operations and pod states
//   - Label selector building for targeted pod queries
//
// The monitoring logic runs periodically based on the configured schedule interval,
// examining pods in specified namespaces and terminating those that exceed the
// maximum allowed lifetime. The package is designed to be robust and handle
// various error conditions gracefully while providing detailed logging.
package monitoring
