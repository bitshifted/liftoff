package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/bitshifted/liftoff/cli"
	"github.com/bitshifted/liftoff/log"
)

var input cli.CLI

func main() {
	ctx := kong.Parse(&input)
	log.Init(debugLoggingEnabled(ctx.Args))
	err := ctx.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Execution failed")
		os.Exit(1)
	}
}

func debugLoggingEnabled(args []string) bool {
	for _, s := range args {
		if s == "--enable-debug" {
			return true
		}
	}
	return false
}
