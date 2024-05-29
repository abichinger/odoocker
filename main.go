package main

import (
	"os"

	"github.com/abichinger/odoocker/flags"
	"github.com/abichinger/odoocker/internal/odoocker"
	"github.com/leaanthony/clir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create a new cli
	cli := clir.NewCli("odoocker", "Run odoo tests inside a docker container", "v1.0.0")

	// Define flags using the Flags struct
	var flags flags.Options
	cli.AddFlags(&flags)

	// Define action for the command
	cli.Action(func() error {

		if len(flags.Modules) == 0 {
			log.Error().Msg("Provide at least one module")
			cli.PrintHelp()
			os.Exit(1)
		}

		if len(flags.TestTags) == 0 {
			log.Error().Msg("Provide at least one test tag")
			cli.PrintHelp()
			os.Exit(1)
		}

		odoo := odoocker.NewOdoocker(&flags)
		return odoo.Run()
	})

	// Run the CLI
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
