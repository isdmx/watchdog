package config

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name         string
		loggingMode  string
		loggingLevel string
		expectError  bool
	}{
		{
			name:         "development mode with info level",
			loggingMode:  "development",
			loggingLevel: "info",
			expectError:  false,
		},
		{
			name:         "production mode with debug level",
			loggingMode:  "production",
			loggingLevel: "debug",
			expectError:  false,
		},
		{
			name:         "invalid log level defaults to info",
			loggingMode:  "production",
			loggingLevel: "invalid_level",
			expectError:  false, // Should not error, just default to info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Logging: LoggingConfig{
					Mode:  tt.loggingMode,
					Level: tt.loggingLevel,
				},
			}

			logger, err := NewLogger(cfg)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && logger == nil {
				t.Error("Expected logger to be created but got nil")
			}

			// Clean up
			if logger != nil {
				_ = logger.Sync()
			}
		})
	}
}

func TestNewSugaredLogger(t *testing.T) {
	// Create a basic logger
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create zap logger: %v", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()

	// Test NewSugaredLogger
	sugaredLogger := NewSugaredLogger(zapLogger)
	if sugaredLogger == nil {
		t.Error("Expected sugared logger to be created but got nil")
	}
}

func TestLoggingModule(t *testing.T) {
	// Test that the module is defined
	if LoggingModule == nil {
		t.Error("LoggingModule should not be nil")
	}
}
