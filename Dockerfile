# Build stage
FROM golang:1.25.3-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o watchdog .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh watchdog

# Set working directory
WORKDIR /home/watchdog

# Copy the binary from builder stage
COPY --from=builder /app/watchdog .

# Change ownership
RUN chown -R watchdog:watchdog /home/watchdog

# Switch to non-root user
USER watchdog

# Expose port for health checks and metrics
EXPOSE 8080

# Run the application
CMD ["./watchdog"]