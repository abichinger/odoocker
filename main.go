package main

import (
	"os"

	"github.com/abichinger/odoocker/internal/odoocker"
	"github.com/leaanthony/clir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create a new cli
	cli := clir.NewCli("odoocker", "Setup and run odoo inside a docker container", "v1.0.0")
	odoocker.AddTestCommand(cli)
	odoocker.AddRunCommand(cli)
	odoocker.AddUpCommand(cli)

	// Run the CLI
	if err := cli.Run(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
