package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/bitshifted/liftoff/cli"
	"github.com/bitshifted/liftoff/log"
)

var input cli.CLI

func main() {
	log.Init()
	ctx := kong.Parse(&input)
	err := ctx.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Execution failed")
		os.Exit(1)
	}
}
