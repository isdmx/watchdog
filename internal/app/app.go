package app

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/isdmx/watchdog/internal/client"
	"github.com/isdmx/watchdog/internal/config"
	"github.com/isdmx/watchdog/internal/logging"
	"github.com/isdmx/watchdog/internal/monitoring"
	"github.com/isdmx/watchdog/internal/server"
)

func NewApplication() *fx.App {
	return fx.New(
		// Configuration module
		fx.Provide(config.NewConfig),

		// Logging module
		fx.Provide(logging.NewLogger),
		fx.Provide(logging.NewSugaredLogger),

		// Kubernetes client module
		fx.Provide(client.NewKubernetesClient),

		// Monitoring module
		fx.Provide(monitoring.NewPodMonitor),

		// HTTP server
		fx.Provide(fx.Annotate(
			server.NewHTTPServer,
			fx.ResultTags(`group:"servers"`),
			fx.As(new(server.Server)),
		)),

		// Watchdog server
		fx.Provide(fx.Annotate(
			server.NewWatchdogServer,
			fx.ResultTags(`group:"servers"`),
			fx.As(new(server.Server)),
		)),

		// Start the application (all servers)
		fx.Invoke(fx.Annotate(
			func([]server.Server) {},
			fx.ParamTags(`group:"servers"`),
		)),

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
	)
}
