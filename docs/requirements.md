# Watchdog Application - Requirement Specification

## Overview

The Watchdog application is a Kubernetes resource management tool designed to monitor and clean up unused resources within a Kubernetes cluster. The primary purpose of this application is to identify and terminate old, unused pods based on configurable criteria to optimize resource utilization and maintain cluster health.

## Main Purpose and Functionality

- **Resource Monitoring**: Continuously monitor designated namespaces for pods that match specified label criteria
- **Resource Cleanup**: Identify and terminate old, unused pods that exceed configured lifetime limits
- **Automated Operations**: Run periodic cleanup operations based on configurable scheduling
- **Configuration-Driven**: Operate based on parameters defined in a configuration file

## Kubernetes Deployment Requirements

- **Deployment Type**: Deployed as a Kubernetes Deployment within the target cluster
- **Namespace**: Can run in any namespace but will monitor designated namespaces for cleanup operations
- **Permissions**: Requires RBAC permissions to list, view, and delete pods in monitored namespaces
- **Service Account**: Requires a dedicated service account with appropriate permissions
- **Resource Limits**: Should define appropriate CPU and memory limits to prevent resource exhaustion

## Resource Cleanup Logic

### Pod Selection Criteria

- Only pods matching specific label keys and values will be considered for monitoring
- Pods must exist in designated namespaces as defined in the configuration
- Pod lifetime will be calculated from its creation timestamp

### Cleanup Decision Process

- Compare pod lifetime against the maximum allowed lifetime defined in the configuration
- Only terminate pods that exceed the lifetime limit and are considered "unused"
- "Unused" pods are determined based on lack of activity or usage metrics (to be defined)

### Execution Schedule

- Cleanup operations run at regular intervals as specified in the configuration file
- Each execution cycle independently evaluates all qualifying pods for potential termination
- Process should be idempotent and safe to run repeatedly

## Configuration File Structure

The application requires a `config.yaml` file at startup with the following structure:

```yaml
watchdog:
  # List of namespaces to monitor
  namespaces:
    - "namespace1"
    - "namespace2"

  # Label selector criteria for identifying target pods
  labelSelectors:
    key1: "value1"
    key2: "value2"

  # Cleanup scheduling in duration format (e.g., "5m", "1h", "24h")
  scheduleInterval: "10m"

  # Maximum pod lifetime before cleanup consideration
  maxPodLifetime: "24h"

  # Additional safety settings
  dryRun: false # If true, only log what would be deleted without taking action

logging:
  # Logging mode: "production" or "development"
  mode: "production"

  # Log level using zap logging levels: "debug", "info", "warn", "error", "dpanic", "panic", "fatal"
  level: "info"
```

## Security Considerations

- Follow security best practices for Kubernetes applications
- Implement minimal required RBAC permissions (least privilege principle)
- Support secure configuration through Kubernetes secrets if sensitive data is involved
- Implement proper authentication and authorization mechanisms

## Performance and Efficiency

- Efficiently handle large numbers of pods without significant performance degradation
- Implement proper resource utilization controls
- Include metrics and monitoring capabilities

## Error Handling and Logging

- Comprehensive error handling for all Kubernetes API interactions
- Detailed logging of cleanup operations and decisions based on logging configuration
- Monitoring and alerting capability for operational visibility
- Graceful handling of transient API errors

## Additional Features

- Health check endpoints for Kubernetes liveness and readiness probes
- Metrics endpoint for Prometheus integration
- Configurable logging settings as defined in the logging section
- Graceful shutdown procedures

## Key Technologies

### Programming Language

- **Go Language**: Implement the entire application using Go for performance, concurrency, and Kubernetes ecosystem compatibility.

### Configuration Management

- **Viper**: Utilize `github.com/spf13/viper` for robust configuration reading, management, and parsing of the config.yaml file.

### Logging Framework

- **Zap Logger**: Use `go.uber.org/zap` for structured, high-performance logging throughout the application.
- **Sugared Logger**: Specifically implement `*zap.SugaredLogger` for easier structured logging with better developer experience.

### Application Framework

- **Fx Framework**: Leverage `go.uber.org/fx` for managing application lifecycle, configuration loading, dependency injection, and graceful shutdown procedures.

### Kubernetes Integration

- **Client-go**: Use `k8s.io/client-go` for all interactions with the Kubernetes cluster, including pod monitoring and management operations.

## Code Quality Requirements

### Simplicity and Readability

- Maintain clean, simple, and highly readable code throughout the application.
- Apply `gofmt` formatting automatically to ensure consistent code style across all files.

### Interface Design

- Design small, focused interfaces that serve specific purposes.
- Follow the Go idiom of returning structs but accepting interfaces as parameters to promote flexibility and testability.

### Function Design

- Avoid unnecessary wrapper functions that don't provide clear abstraction or functionality benefits.
- Only introduce wrapper functions when they serve a clear purpose in abstraction or enhance functionality.

### Testing Strategy

- Implement comprehensive tests for all functions, components, and features.
- Place test files in `_test.go` files alongside the code they test for better organization and maintainability.
- Ensure tests cover both happy path and error scenarios for all critical functionality.
