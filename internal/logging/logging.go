package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/isdmx/watchdog/internal/config"
)

// NewLogger creates a new zap logger based on the configuration
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if cfg.Logging.Mode == "development" {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	// Set the log level based on configuration
	level, err := zapcore.ParseLevel(cfg.Logging.Level)
	if err != nil {
		// Default to info level if parsing fails
		level = zap.InfoLevel
	}
	loggerConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewSugaredLogger wraps the logger with SugaredLogger
func NewSugaredLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
