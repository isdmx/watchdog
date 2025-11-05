package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// MockServer implements the Server interface for testing
type MockServer struct {
	startFunc    func(ctx context.Context) error
	shutdownFunc func(ctx context.Context) error
}

func (m *MockServer) Start(ctx context.Context) error {
	if m.startFunc != nil {
		return m.startFunc(ctx)
	}
	return nil
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc(ctx)
	}
	return nil
}

func TestServerInterface(t *testing.T) {
	t.Run("interface is properly defined", func(t *testing.T) {
		// This test just verifies that the interface exists and can be implemented
		var s Server
		require.Nil(t, s)

		// Create a mock server to verify interface compatibility
		mockServer := &MockServer{}
		s = mockServer

		require.NotNil(t, s)
	})
}
