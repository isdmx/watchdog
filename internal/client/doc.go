// Package client provides Kubernetes client functionality for interacting with the Kubernetes API.
//
// This package handles the creation and configuration of Kubernetes clients, including
// support for both in-cluster and out-of-cluster configurations. It abstracts the
// complexities of Kubernetes client setup and provides a clean interface for other
// packages to interact with the Kubernetes API.
//
// The package includes:
//   - Automatic detection and use of in-cluster configuration when running inside Kubernetes
//   - Fallback to kubeconfig file when running outside the cluster (e.g., for local development)
//   - Proper error handling and logging for client creation failures
//   - Integration with the application's logging system
//
// The NewKubernetesClient function follows the dependency injection pattern used
// throughout the application and returns a configured Kubernetes client interface
// that can be used by other components to perform operations on Kubernetes resources.
package client
