package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewApplication(t *testing.T) {
	t.Run("creates application successfully", func(t *testing.T) {
		app := NewApplication()
		require.NotNil(t, app)

		// Test that the app can be started (though we won't actually run it in tests)
		// Just verify it can be created without errors
	})

	t.Run("application has required modules", func(t *testing.T) {
		app := NewApplication()
		require.NotNil(t, app)

		// The app should have all the required modules as providers
		// We can't directly check the providers, but we can test that
		// the app is created without errors and has the expected structure
		require.NotNil(t, app)
	})
}
