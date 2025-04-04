package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/sch8ill/propmon/api"
	"github.com/sch8ill/propmon/broker"
	"github.com/sch8ill/propmon/config"
	"github.com/sch8ill/propmon/proposal"
	"github.com/sch8ill/propmon/proposal/expiration"
	"github.com/sch8ill/propmon/quality"
)

func main() {
	createLogger()

	app := createApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Error encountered")
	}
}

func monitorProposals(ctx *cli.Context) error {
	config.SetConfig(ctx)

	r := proposal.NewProposalRepository(config.ProposalLifetime)

	listener := broker.NewListener(config.BrokerAddress, r)
	if err := listener.Listen(); err != nil {
		return fmt.Errorf("failed to start broker listener: %w", err)
	}
	defer listener.Shutdown()

	expirationService := expiration.NewExpirationService(r, config.ExpirationJobInterval)
	expirationService.Start()
	defer expirationService.Stop()

	qualityOracle := quality.NewOracle(config.QualityOracle)
	qualityService := quality.NewQualityService(qualityOracle, r, config.QualityUpdateInterval, config.ProposalLifetime)
	qualityService.Start()
	defer qualityService.Stop()

	apiServer := api.New(config.MetricsAddress, r)
	if err := apiServer.Run(); err != nil {
		return fmt.Errorf("api server: %w", err)
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
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}
	log.Logger = log.Output(consoleWriter).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}
