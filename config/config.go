package config

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	DefaultBrokerAddress         string = "nats://broker.mysterium.network:4222"
	DefaultMetricsAddress               = ":9500"
	DefaultProposalLifetime             = 3*time.Minute + 10*time.Second
	DefaultExpirationJobInterval        = 20 * time.Second
	DefaultQualityOracle                = "https://quality.mysterium.network"
	DefaultQualityUpdateInterval        = 30 * time.Minute

	BrokerAddressFlag         = "broker-address"
	MetricsAddressFlag        = "metrics-address"
	ProposalLifetimeFlag      = "proposal-lifetime"
	ExpirationJobIntervalFlag = "expiration-job-delay"
	QualityOracleFlag         = "quality-oracle"
	QualityUpdateIntervalFlag = "quality-update-interval"
)

var (
	BrokerAddress         string
	MetricsAddress        string
	ProposalLifetime      time.Duration
	ExpirationJobInterval time.Duration
	QualityOracle         string
	QualityUpdateInterval time.Duration
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
			Name:  ExpirationJobIntervalFlag,
			Usage: "delay between expiration job runs",
			Value: DefaultExpirationJobInterval,
		},
		&cli.StringFlag{
			Name:  MetricsAddressFlag,
			Usage: "address the prometheus metrics exporter listens on",
			Value: DefaultMetricsAddress,
		},
		&cli.StringFlag{
			Name:  QualityOracleFlag,
			Usage: "url of the quality oracle",
			Value: DefaultQualityOracle,
		},
		&cli.DurationFlag{
			Name:  QualityUpdateIntervalFlag,
			Usage: "interval between quality data updates",
			Value: DefaultQualityUpdateInterval,
		},
	}
}

func SetConfig(ctx *cli.Context) {
	BrokerAddress = ctx.String(BrokerAddressFlag)
	MetricsAddress = ctx.String(MetricsAddressFlag)
	ProposalLifetime = ctx.Duration(ProposalLifetimeFlag)
	ExpirationJobInterval = ctx.Duration(ExpirationJobIntervalFlag)
	QualityOracle = ctx.String(QualityOracleFlag)
	QualityUpdateInterval = ctx.Duration(QualityUpdateIntervalFlag)
}
