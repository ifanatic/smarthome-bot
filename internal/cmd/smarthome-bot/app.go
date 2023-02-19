package smarthomebot

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func AppModule(configFile string) fx.Option {
	return fx.Options(
		LoggerModule,
		TelegramModule,
		fx.Decorate(func(logger *zap.Logger) *zap.Logger {
			logger.Named("smarthome-bot")
			return logger
		}),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(*TelegramBot) {}),
	)
}
