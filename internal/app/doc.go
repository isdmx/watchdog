// Package app contains the main application setup and dependency injection configuration.
//
// This package uses the Uber fx framework to manage the application lifecycle and
// dependencies. It orchestrates the initialization of all major components including
// configuration, logging, Kubernetes client, monitoring, and HTTP server services.
//
// The NewApplication function serves as the main entry point that wires together
// all the individual components into a cohesive application using fx's dependency
// injection capabilities. The application is structured as a modular system where
// each major functionality (config, logging, client, monitoring, server) is handled
// by its own package and injected where needed.
//
// Key features of this package:
//   - Centralized dependency injection using Uber fx
//   - Modular architecture with clear separation of concerns
//   - Lifecycle management for all application components
//   - Integration with fx logging for dependency injection events
package app
