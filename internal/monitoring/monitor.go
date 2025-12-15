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
	ager := NewAgerFromConfig(&pm.config.Watchdog, pm.logger)

	for _, namespace := range pm.config.Watchdog.Namespaces {
		logger_namespace := pm.logger.WithLazy("namespace", namespace)
		logger_namespace.Debugf("Processing namespace")

		// List pods in the namespace with the specified labels
		pods, err := pm.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			logger_namespace.Errorw("Failed to list pods", "error", err)
			continue
		}

		logger_namespace.Debugf("Found %d pods in namespace with matching labels", len(pods.Items))
		PodsExaminedTotal.Add(float64(len(pods.Items)))

		// Filter and terminate old pods
		for i := range pods.Items {
			pod := &pods.Items[i]
			logger_pod := logger_namespace.WithLazy("pod", pod.Name)

			isOld, err := ager.IsOld(pod)
			if err != nil {
				logger_pod.Warnw("Unable to calculate pod age", "pod", pod.Name, "err", err)
				continue
			}
			if !isOld {
				continue
			}

			if pm.config.Watchdog.DryRun {
				logger_pod.Infow("DRY RUN: Would terminate pod")
				PodsTerminatedTotal.WithLabelValues(namespace, "true").Inc()
			} else {
				// Terminate the pod
				err := pm.terminatePod(namespace, pod.Name)
				if err != nil {
					logger_pod.Errorw("Failed to terminate pod", "error", err)
				} else {
					logger_pod.Infow("Successfully terminated pod")
					PodsTerminatedTotal.WithLabelValues(namespace, "false").Inc()
					PodsTerminatedByAgeTotal.WithLabelValues(namespace).Inc()
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
