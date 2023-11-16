package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/sch8ill/propmon/broker"
	"github.com/sch8ill/propmon/config"
	"github.com/sch8ill/propmon/metrics"
	"github.com/sch8ill/propmon/proposal"
	"github.com/sch8ill/propmon/proposal/expiration"
)

func main() {
	createLogger()

	app := createApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("error encountered")
	}
}

func monitorProposals(ctx *cli.Context) error {
	config.SetConfig(ctx)

	r := proposal.NewProposalRepository(config.ProposalLifetime)
	expirationService := expiration.NewExpirationService(r, config.ExpirationJobDelay)
	listener := broker.NewListener(config.BrokerAddress, r)
	expirationService.Start()
	defer expirationService.Stop()

	if err := listener.Listen(); err != nil {
		return fmt.Errorf("failed to start broker listener: %w", err)
	}
	defer listener.Shutdown()

	if err := metrics.Listen(); err != nil {
		return fmt.Errorf("failed to start prometheus exporter: %w", err)
	}

	return nil
}

func createApp() *cli.App {
	return &cli.App{
		Name:      "propmon",
		Usage:     "monitor mysterium network node service proposals",
		Copyright: "Copyright (c) 2023 Sch8ill",
		Action:    monitorProposals,
		Flags:     config.DeclareFlags(),
	}
}

func createLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}
	log.Logger = log.Output(consoleWriter).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}
