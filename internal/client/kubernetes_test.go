package client

import (
	"testing"
)

func TestKubernetesModule(t *testing.T) {
	// Simple test to ensure the module is defined
	if KubernetesModule == nil {
		t.Error("KubernetesModule should not be nil")
	}
}
