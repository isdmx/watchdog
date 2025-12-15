package monitoring

import (
	"strconv"
	"time"

	"github.com/isdmx/watchdog/internal/config"

	"go.uber.org/zap"

	k8type "k8s.io/api/core/v1"
)

var _ Ager = (*CreationAger)(nil)
var _ Ager = (*LabeledAger)(nil)

type Ager interface {
	IsOld(*k8type.Pod) (bool, error)
}

type CreationAger struct {
	maxPodLifetime time.Duration
	logger         *zap.SugaredLogger
}

func (a *CreationAger) IsOld(pod *k8type.Pod) (bool, error) {
	logger := a.logger.With("pod", pod.Name, "namespace", pod.Namespace)

	age := time.Since(pod.CreationTimestamp.Time)
	logger.Debugf("Pod %s age: %v, max age: %v", pod.Name, age, a.maxPodLifetime)

	if age <= a.maxPodLifetime {
		return false, nil
	}

	logger.Infow("Pod exceeds maximum lifetime",
		"age", age,
		"maxAge", a.maxPodLifetime,
	)
	return true, nil
}

func NewCreationAger(maxPodLifetime time.Duration, logger *zap.SugaredLogger) *CreationAger {
	return &CreationAger{
		maxPodLifetime: maxPodLifetime,
		logger:         logger.WithLazy("ager", "creation"),
	}
}

type LabeledAger struct {
	labelKillTime  string
	maxPodLifetime time.Duration
	logger         *zap.SugaredLogger
}

func (a *LabeledAger) IsOld(pod *k8type.Pod) (bool, error) {
	logger := a.logger.With("pod", pod.Name, "namespace", pod.Namespace)

	age := time.Since(pod.CreationTimestamp.Time)
	logger.Debugf("Pod %s age: %v, max age: %v, labels: %v", pod.Name, age, a.maxPodLifetime, pod.Labels)

	if a.maxPodLifetime <= age {
		// Using pod.CreationTimestamp is not desired but good as fallback path
		logger.Warnw("Terminating pod by creation time", "age", age)
		return true, nil
	}

	killTimeRaw, exists := pod.Labels[a.labelKillTime]
	if !exists {
		logger.Warnw("No ttl label in pod")
		return false, nil
	}

	killTime, err := strconv.ParseFloat(killTimeRaw, 64)
	if err != nil {
		return false, err
	}
	if killTime <= float64(time.Now().Unix()) {
		logger.Infow("Killing pod by TTL", "kill_time", killTime)
		return true, nil
	}
	return false, nil
}

func NewLabeledAger(labelKillTime string, maxPodLifetime time.Duration, logger *zap.SugaredLogger) *LabeledAger {
	return &LabeledAger{
		labelKillTime:  labelKillTime,
		maxPodLifetime: maxPodLifetime,
		logger:         logger.WithLazy("ager", "creation+label", "label", labelKillTime),
	}
}

func NewAgerFromConfig(cfg *config.WatchdogConfig, logger *zap.SugaredLogger) Ager {
	if cfg.TtlLabel == "" {
		return NewCreationAger(cfg.MaxPodLifetime, logger)
	}
	return NewLabeledAger(cfg.TtlLabel, cfg.MaxPodLifetime, logger)
}
