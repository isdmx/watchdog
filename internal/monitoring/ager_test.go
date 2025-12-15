package monitoring

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/isdmx/watchdog/internal/config"
	"go.uber.org/zap"
	k8type "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreationAger(t *testing.T) {
	ager := NewCreationAger(time.Hour*2, zap.NewNop().Sugar())
	tests := []struct {
		name        string
		pod         k8type.Pod
		expectedErr error
		expected    bool
		ager        *CreationAger
	}{
		{
			name: "simple",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-4 * time.Hour)},
				},
			},
			expected: true,
			ager:     ager,
		},
		{
			name: "now",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: false,
			ager:     ager,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ager.IsOld(&tt.pod)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			}
			require.Equal(t, tt.expected, res)

		})
	}
}

func TestLabeledAger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ager := NewLabeledAger("sandbox.kill_time", time.Hour*2, logger.Sugar()) //zap.NewNop().Sugar())
	tests := []struct {
		name        string
		pod         k8type.Pod
		expectedErr error
		expected    bool
		ager        *LabeledAger
	}{
		{
			name: "creation-simple",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-4 * time.Hour)},
				},
			},
			expected: true,
			ager:     ager,
		},
		{
			name: "creation-now",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: false,
			ager:     ager,
		},
		{
			name: "no-labels",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: false,
			ager:     ager,
		},
		{
			name: "labeled-recent",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					Labels:            map[string]string{"sandbox.kill_time": fmt.Sprintf("%v", time.Now().Add(time.Hour).Unix())},
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: false,
			ager:     ager,
		},
		{
			name: "labeled-old",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					Labels:            map[string]string{"sandbox.kill_time": fmt.Sprintf("%v", time.Now().Add(-time.Minute).Unix())},
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: true,
			ager:     ager,
		},

		{
			name: "labeled-old-float",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					Labels:            map[string]string{"sandbox.kill_time": fmt.Sprintf("%v", 0.05123+float64(time.Now().Add(-time.Minute).Unix()))},
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expected: true,
			ager:     ager,
		},
		{
			name: "labeled-invalid_time",
			pod: k8type.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "old-pod",
					Namespace:         "default",
					Labels:            map[string]string{"sandbox.kill_time": "axaxa"},
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expectedErr: strconv.ErrSyntax,
			expected:    false,
			ager:        ager,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ager.IsOld(&tt.pod)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, res)

		})
	}
}

func TestNewAgerFromConfig(t *testing.T) {
	emptyRes := NewAgerFromConfig(&config.WatchdogConfig{}, zap.NewNop().Sugar())
	require.IsType(t, &CreationAger{}, emptyRes)

	configWithTtl := &config.WatchdogConfig{}
	configWithTtl.TtlLabel = "dasdfasd"
	require.IsType(t, &LabeledAger{}, NewAgerFromConfig(configWithTtl, zap.NewNop().Sugar()))
}
