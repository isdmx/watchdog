package client

import (
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// NewKubernetesClient creates a new Kubernetes client
func NewKubernetesClient(logger *zap.SugaredLogger) (kubernetes.Interface, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Try in-cluster config first (for when running inside Kubernetes)
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Warnw("Failed to get in-cluster config", "kubeconfig", kubeconfig, "error", err)

		// Fall back to kubeconfig file (for local development)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Errorw("Failed to create Kubernetes client", "error", err)

		return nil, err
	}

	return clientset, nil
}
