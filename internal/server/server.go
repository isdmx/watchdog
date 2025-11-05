package server

import (
	"context"
)

// Server interface defines the contract for all servers in the application
type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
