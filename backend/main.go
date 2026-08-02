// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bitshifted/fundslock-be/cli"
	"bitshifted/fundslock-be/log"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

var input cli.CLI

func main() {
	// Only load .env if running locally; will not crash if missing in production
	_ = godotenv.Load()
	ctx := kong.Parse(&input)
	log.Init(debugLoggingEnabled(ctx.Args))
	log.Logger.Info().Msg("Starting server...")
	log.Logger.Debug().Msg("Debug logging enabled")
	err := ctx.Run()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to run command")
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
