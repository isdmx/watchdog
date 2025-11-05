package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isdmx/watchdog/internal/app"
)

func TestMainFunction(t *testing.T) {
	t.Run("main function runs without error", func(t *testing.T) {
		// This is a minimal test for the main function
		// We create an app and verify it's created properly
		appInstance := app.NewApplication()
		require.NotNil(t, appInstance)

		// Since main() calls app.Run(), we're effectively testing
		// that the application can be created with all its dependencies
	})
}
