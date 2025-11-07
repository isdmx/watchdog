// Package logging provides structured logging functionality for the application.
//
// This package sets up and configures the Zap logging library based on the
// application's configuration. It provides both structured and sugared logging
// interfaces to suit different logging needs throughout the application.
//
// The logging configuration supports:
//   - Development and production logging modes
//   - Configurable log levels (debug, info, warn, error, etc.)
//   - Structured logging with key-value pairs for better log analysis
//   - ISO8601 timestamp formatting
//   - Integration with the application's dependency injection system
//
// The package provides two main functions:
//   - NewLogger: Creates a configured Zap logger instance based on the application config
//   - NewSugaredLogger: Wraps the structured logger with a sugared logger for easier use
//
// The logging setup is designed to be consistent across all application components
// while providing the flexibility to adjust logging behavior based on configuration.
package logging
