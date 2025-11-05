package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// PodsTerminatedTotal counts the total number of pods terminated
	PodsTerminatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchdog_pods_terminated_total",
			Help: "Total number of pods terminated by the watchdog",
		},
		[]string{"namespace", "dry_run"},
	)

	// MonitoringDuration tracks how long monitoring runs take
	MonitoringDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "watchdog_monitoring_duration_seconds",
			Help:    "Time spent running monitoring checks",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
	)

	// PodsExaminedTotal counts the total number of pods examined
	PodsExaminedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "watchdog_pods_examined_total",
			Help: "Total number of pods examined by the watchdog",
		},
	)

	// PodsTerminatedByAgeTotal counts pods terminated due to age
	PodsTerminatedByAgeTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchdog_pods_terminated_by_age_total",
			Help: "Total number of pods terminated due to age limits",
		},
		[]string{"namespace"},
	)
)
