# Watchdog Makefile

# Build the application
.PHONY: build
build:
	go build -o bin/watchdog ./cmd/watchdog

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage and show in terminal
.PHONY: coverage
coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run tests with coverage and generate HTML report
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

# Run the application locally
.PHONY: run
run:
	go run ./cmd/watchdog

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t watchdog:latest .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --rm watchdog:latest

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Lint code
.PHONY: lint
lint:
	golangci-lint run ./...

# Install golangci-lint
.PHONY: install-lint
install-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

# Run go fmt
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet
.PHONY: vet
vet:
	go vet ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f bin/watchdog
	rm -f coverage.out
	rm -f coverage.html

# Deploy to Kubernetes
.PHONY: deploy
deploy:
	kubectl apply -f deployments/k8s/configmap.yaml
	kubectl apply -f deployments/k8s/rbac.yaml
	kubectl apply -f deployments/k8s/deployment.yaml

# Remove deployment from Kubernetes
.PHONY: undeploy
undeploy:
	kubectl delete -f deployments/k8s/deployment.yaml
	kubectl delete -f deployments/k8s/rbac.yaml
	kubectl delete -f deployments/k8s/configmap.yaml

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage and show in terminal"
	@echo "  test-coverage - Run tests with coverage and generate HTML report"
	@echo "  run           - Run the application locally"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  clean         - Clean build artifacts"
	@echo "  deploy        - Deploy to Kubernetes"
	@echo "  undeploy      - Remove deployment from Kubernetes"
	@echo "  help          - Show this help message"
