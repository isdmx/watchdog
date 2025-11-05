# Watchdog

Watchdog is a Kubernetes resource management tool that monitors and cleans up unused resources within a Kubernetes cluster. The primary purpose of this application is to identify and terminate old, unused pods based on configurable criteria to optimize resource utilization and maintain cluster health.

## Features

- **Resource Monitoring**: Continuously monitors designated namespaces for pods that match specified label criteria
- **Resource Cleanup**: Identifies and terminates old, unused pods that exceed configured lifetime limits
- **Automated Operations**: Runs periodic cleanup operations based on configurable scheduling
- **Configuration-Driven**: Operates based on parameters defined in a configuration file
- **Health Checks**: Provides health check endpoints for Kubernetes liveness and readiness probes
- **Metrics**: Exposes Prometheus metrics for monitoring
- **Dry Run Mode**: Allows testing cleanup operations without actually terminating pods

## Requirements

- Go 1.25+
- Kubernetes cluster access

## Configuration

The application requires a `config.yaml` file with the following structure:

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
  dryRun: false  # If true, only log what would be deleted without taking action

logging:
  # Logging mode: "production" or "development"
  mode: "production"

  # Log level using zap logging levels: "debug", "info", "warn", "error", "dpanic", "panic", "fatal"
  level: "info"
```

## Building

To build the application:

```bash
go build -o watchdog ./cmd/watchdog
```

## Running Locally

To run the application locally:

```bash
# Make sure you have a valid kubeconfig file
go run main.go
```

## Kubernetes Deployment

The application is designed to run as a deployment in Kubernetes. See the `k8s/` directory for deployment manifests.

To deploy to Kubernetes:

```bash
kubectl apply -f deployments/k8s/
```

## Endpoints

- `/healthz` - Health check endpoint
- `/readyz` - Readiness check endpoint
- `/metrics` - Prometheus metrics endpoint

## Development

Run tests:

```bash
go test ./...
```

Format code:

```bash
go fmt ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

### Setting up Development Tools

#### Installing pre-commit

To install pre-commit, you can use the following methods:

- **Using pip**:

  ```bash
  pip install pre-commit
  ```

- **Using Homebrew (macOS)**:

  ```bash
  brew install pre-commit
  ```

For more information, visit [pre-commit's official website](https://pre-commit.com/).

#### Installing golangci-lint

To install golangci-lint, you can use the following methods:

- **Using Go install**:

  ```bash
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
  ```

- **Using Homebrew (macOS)**:

  ```bash
  brew install golangci-lint
  ```

For more information, visit [golangci-lint's official website](https://golangci-lint.run/).

## License

MIT
