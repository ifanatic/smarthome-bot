package smarthomebot

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ifanatic/smarthome-bot/config"
)

type TelegramConfig struct {
	Token         string
	UpdateTimeout int
	Debug         bool
}

type TelegramBot struct {
	logger *zap.Logger

	config *TelegramConfig
	client *botapi.BotAPI

	shutdownCh chan interface{}
}

func NewTelegramBot(logger *zap.Logger, config *TelegramConfig) (*TelegramBot, error) {
	client, err := botapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	client.Debug = config.Debug

	return &TelegramBot{
		logger:     logger,
		config:     config,
		client:     client,
		shutdownCh: make(chan interface{}),
	}, nil
}

func (b *TelegramBot) Start(ctx context.Context) error {
	b.logger.Info("Starting Telegram Bot")

	go b.HandleUpdates()

	return nil
}

func (b *TelegramBot) Stop(ctx context.Context) error {
	b.logger.Info("Stopping Telegram Bot")
	close(b.shutdownCh)
	return nil
}

func (b *TelegramBot) HandleUpdates() {
	u := botapi.NewUpdate(0)
	u.Timeout = b.config.UpdateTimeout

	updates := b.client.GetUpdatesChan(u)

	b.logger.Info("Listen for updates")

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				b.logger.Error("Updates channel unexpectedly closed, stopping")
				return
			}
			b.handleUpdate(&update)
		case <-b.shutdownCh:
			b.logger.Info("Shutdown request received, stopping")
			return
		}
	}
}

func (b *TelegramBot) handleUpdate(u *botapi.Update) {
	if u.Message == nil || !u.Message.IsCommand() {
		b.logger.Debug("Not a command")
	}

	msg := botapi.NewMessage(u.Message.Chat.ID, "")

	switch u.Message.Command() {
	case "status":
		msg.Text = "Command \"status\" received"
	default:
		msg.Text = "I don't known that command"
	}

	if _, err := b.client.Send(msg); err != nil {
		b.logger.Sugar().Errorw("Failed to send a message",
			"error", err,
		)
	}
}

type TelegramModuleParams struct {
	fx.In

	Logger *zap.Logger
	Config *config.Config
}

type TelegramModuleResult struct {
	fx.Out

	TelegramBot *TelegramBot
}

var TelegramModule = fx.Provide(func(lc fx.Lifecycle, params TelegramModuleParams) (TelegramModuleResult, error) {
	config := TelegramConfig{
		Token:         params.Config.Telegram.Token,
		UpdateTimeout: params.Config.Telegram.UpdateTimeout,
		Debug:         params.Config.Telegram.Debug,
	}

	bot, err := NewTelegramBot(params.Logger, &config)
	if err != nil {
		return TelegramModuleResult{}, err
	}

	lc.Append(fx.Hook{
		OnStart: bot.Start,
		OnStop:  bot.Stop,
	})
	return TelegramModuleResult{
		TelegramBot: bot,
	}, nil
})
