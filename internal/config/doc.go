// Package config handles application configuration management.
//
// This package provides structured configuration for the watchdog application,
// using the Viper library for flexible configuration loading. It supports
// configuration from multiple sources including YAML files and environment
// variables, with sensible default values for all required parameters.
//
// The configuration is organized into logical sections:
//   - HttpConfig: HTTP server settings (address, timeouts)
//   - WatchdogConfig: Pod monitoring settings (namespaces, selectors, intervals, limits)
//   - LoggingConfig: Logging settings (mode, level)
//
// The package provides:
//   - Default configuration values for all settings
//   - Support for both development and production configuration modes
//   - Type-safe configuration access through structured configuration structs
//   - Integration with the application's dependency injection system
//
// Configuration can be loaded from:
//   - config.yaml in the current directory
//   - config/config.yaml in the config subdirectory
//   - Environment variables (with proper prefixes)
package config
