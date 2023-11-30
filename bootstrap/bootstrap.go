package bootstrap

import (
	config "erp/config"
	"erp/internal/api/controllers"
	"erp/internal/api/middlewares"
	"erp/internal/api/route"
	"erp/internal/infrastructure"
	"erp/internal/lib"
	"erp/internal/repository"
	"erp/internal/services"
	"erp/internal/utils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func inject() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		//fx.NopLogger,
		fx.Provide(
			config.NewConfig,
			utils.NewTimeoutContext,
		),
		route.Module,
		lib.Module,
		repository.Module,
		service.Module,
		controller.Module,
		middlewares.Module,
		infrastructure.Module,
	)
}

func Run() {
	fx.New(inject()).Run()
}
