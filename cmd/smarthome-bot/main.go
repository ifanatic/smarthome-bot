package main

import (
	"fmt"
	"os"

	"github.com/ifanatic/smarthome-bot/internal/cmd/smarthome-bot"
)

func main() {
	app := smarthomebot.BuildCLI()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start application: %s", err)
		os.Exit(1)
	}
}
