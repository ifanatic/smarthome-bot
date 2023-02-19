package smarthomebot

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var LoggerModule = fx.Provide(func(lc fx.Lifecycle) (*zap.Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			l.Sync()
			return nil
		},
	})
	return l, nil
})
