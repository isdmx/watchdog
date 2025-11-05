package monitoring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	t.Run("Metrics are properly initialized", func(t *testing.T) {
		// Verify that all metrics are properly initialized
		require.NotNil(t, PodsTerminatedTotal)
		require.NotNil(t, MonitoringDuration)
		require.NotNil(t, PodsExaminedTotal)
		require.NotNil(t, PodsTerminatedByAgeTotal)

		// Test that we can use the metrics without errors
		labels := map[string]string{"namespace": "test", "dry_run": "false"}
		PodsTerminatedTotal.With(labels).Inc()

		// Observe a duration
		MonitoringDuration.Observe(0.5) // 0.5 seconds

		// Increment the counter
		PodsExaminedTotal.Inc()

		// Increment the counter with namespace label
		ageLabels := map[string]string{"namespace": "test"}
		PodsTerminatedByAgeTotal.With(ageLabels).Inc()
	})
}
