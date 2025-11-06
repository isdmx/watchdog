package monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"
)

// PodMonitor handles pod monitoring and cleanup operations
type PodMonitor struct {
	clientset kubernetes.Interface
	config    *config.Config
	logger    *zap.SugaredLogger
}

// NewPodMonitor creates a new pod monitor
func NewPodMonitor(clientset kubernetes.Interface, cfg *config.Config, logger *zap.SugaredLogger) *PodMonitor {
	return &PodMonitor{
		clientset: clientset,
		config:    cfg,
		logger:    logger.Named("PodMonitor"),
	}
}

// MonitorAndCleanup performs the monitoring and cleanup operation
func (pm *PodMonitor) MonitorAndCleanup() error {
	pm.logger.Info("Starting pod monitoring and cleanup")

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		MonitoringDuration.Observe(duration.Seconds())
		pm.logger.Debugf("Monitoring completed in %v", duration)
	}()

	// Create label selector from config
	labelSelector := buildLabelSelector(pm.config.Watchdog.LabelSelectors)

	for _, namespace := range pm.config.Watchdog.Namespaces {
		pm.logger.Debugf("Processing namespace: %s", namespace)

		// List pods in the namespace with the specified labels
		pods, err := pm.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			pm.logger.Errorw("Failed to list pods", "namespace", namespace, "error", err)
			continue
		}

		pm.logger.Debugf("Found %d pods in namespace %s with matching labels", len(pods.Items), namespace)
		PodsExaminedTotal.Add(float64(len(pods.Items)))

		// Filter and terminate old pods
		for i := range pods.Items {
			pod := &pods.Items[i] // Use pointer to avoid copying
			age := time.Since(pod.CreationTimestamp.Time)

			pm.logger.Debugf("Pod %s age: %v, max age: %v", pod.Name, age, pm.config.Watchdog.MaxPodLifetime)

			// Check if the pod exceeds the maximum lifetime
			if age > pm.config.Watchdog.MaxPodLifetime {
				pm.logger.Infow("Pod exceeds maximum lifetime",
					"pod", pod.Name,
					"namespace", namespace,
					"age", age,
					"maxAge", pm.config.Watchdog.MaxPodLifetime)

				if pm.config.Watchdog.DryRun {
					pm.logger.Infow("DRY RUN: Would terminate pod", "pod", pod.Name, "namespace", namespace)
					PodsTerminatedTotal.WithLabelValues(namespace, "true").Inc()
				} else {
					// Terminate the pod
					err := pm.terminatePod(namespace, pod.Name)
					if err != nil {
						pm.logger.Errorw("Failed to terminate pod", "pod", pod.Name, "namespace", namespace, "error", err)
					} else {
						pm.logger.Infow("Successfully terminated pod", "pod", pod.Name, "namespace", namespace)
						PodsTerminatedTotal.WithLabelValues(namespace, "false").Inc()
						PodsTerminatedByAgeTotal.WithLabelValues(namespace).Inc()
					}
				}
			}
		}
	}

	return nil
}

// buildLabelSelector creates a label selector string from a map
func buildLabelSelector(labels map[string]string) string {
	selectorParts := make([]string, 0, len(labels))
	for key, value := range labels {
		selectorParts = append(selectorParts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(selectorParts, ",")
}

// terminatePod terminates a pod in the specified namespace
func (pm *PodMonitor) terminatePod(namespace, podName string) error {
	err := pm.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	return err
}
