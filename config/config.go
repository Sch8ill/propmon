package config

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	DefaultBrokerAddress      string = "nats://broker.mysterium.network:4222"
	DefaultMetricsAddress            = ":9500"
	DefaultProposalLifetime          = 3*time.Minute + 10*time.Second
	DefaultExpirationJobDelay        = 20 * time.Second

	BrokerAddressFlag      = "broker-address"
	MetricsAddressFlag     = "metrics-address"
	ProposalLifetimeFlag   = "proposal-lifetime"
	ExpirationJobDelayFlag = "expiration-job-delay"
)

var (
	BrokerAddress      string
	MetricsAddress     string
	ProposalLifetime   time.Duration
	ExpirationJobDelay time.Duration
)

func DeclareFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  BrokerAddressFlag,
			Usage: "broker address to listen for proposals",
			Value: DefaultBrokerAddress,
		},
		&cli.DurationFlag{
			Name:  ProposalLifetimeFlag,
			Usage: "lifetime of a proposal until it expires if not renewed",
			Value: DefaultProposalLifetime,
		},
		&cli.DurationFlag{
			Name:  ExpirationJobDelayFlag,
			Usage: "delay between expiration job runs",
			Value: DefaultExpirationJobDelay,
		},
		&cli.StringFlag{
			Name:  MetricsAddressFlag,
			Usage: "address the prometheus metrics exporter listens on",
			Value: DefaultMetricsAddress,
		},
	}
}

func SetConfig(ctx *cli.Context) {
	BrokerAddress = ctx.String(BrokerAddressFlag)
	MetricsAddress = ctx.String(MetricsAddressFlag)
	ProposalLifetime = ctx.Duration(ProposalLifetimeFlag)
	ExpirationJobDelay = ctx.Duration(ExpirationJobDelayFlag)
}
