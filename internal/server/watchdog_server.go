package server

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/config"
	"github.com/isdmx/watchdog/internal/monitoring"
)

var _ Server = (*WatchdogServer)(nil)

type WatchdogServer struct {
	pm          *monitoring.PodMonitor
	logger      *zap.SugaredLogger
	config      *config.Config
	stopChannel chan struct{}
}

func NewWatchdogServer(lc fx.Lifecycle, pm *monitoring.PodMonitor, logger *zap.SugaredLogger, cfg *config.Config) *WatchdogServer {
	wd := &WatchdogServer{
		pm:          pm,
		logger:      logger,
		config:      cfg,
		stopChannel: make(chan struct{}),
	}

	lc.Append(fx.Hook{
		OnStart: wd.Start,
		OnStop:  wd.Shutdown,
	})

	return wd
}

// StartMonitoring starts the monitoring process
func (wd *WatchdogServer) Start(_ context.Context) error {
	wd.logger.Info("Starting periodic monitoring", "interval", wd.config.Watchdog.ScheduleInterval)

	go func() {
		// Start periodic monitoring
		ticker := time.NewTicker(wd.config.Watchdog.ScheduleInterval)

		for {
			select {
			case <-ticker.C:
				wd.logger.Info("Starting scheduled monitoring check")
				err := wd.pm.MonitorAndCleanup()
				if err != nil {
					wd.logger.Error("Scheduled monitoring run failed", "error", err)
				}
			case <-wd.stopChannel:
				wd.logger.Info("Stopping monitoring")
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

// Stop stops the monitoring process
func (wd *WatchdogServer) Shutdown(_ context.Context) error {
	close(wd.stopChannel)
	wd.logger.Info("Monitoring stopped")
	return nil
}
