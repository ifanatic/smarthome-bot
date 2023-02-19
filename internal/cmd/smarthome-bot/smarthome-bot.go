package smarthomebot

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"github.com/ifanatic/smarthome-bot/config"
)

func BuildCLI() *cli.App {
	app := &cli.App{
		Name: "smarthome-bot",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "relative or absolule path to configuration file (yaml)",
				Value:   "./config/config.yaml",
				Aliases: []string{"c"},
				EnvVars: []string{"SMARTHOME_CONFIG"},
			},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "start smarthome-bot",
			Action: func(cCtx *cli.Context) error {
				return startBot(cCtx.String("config"))
			},
		},
	}

	return app
}

func startBot(configFile string) error {
	cfg, err := config.LoadConfigFromFile(configFile)
	if err != nil {
		return err
	}

	fx.New(
		fx.Provide(func() *config.Config {
			return cfg
		}),
		AppModule(configFile),
	).Run()

	return nil
}
