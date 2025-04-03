package metrics

import (
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sch8ill/propmon/proposal"
)

// custom registry to discard default go metrics
var Registry = prometheus.NewRegistry()

type providerLabel struct {
	Country  string
	NodeType string
}

var proposalRegistered = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "propmon_proposal_registered",
	Help: "Service Proposal registered",
})

var proposalPing = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "propmon_proposal_ping",
	Help: "Service Proposal ping",
})

var proposalUnregistered = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "propmon_proposal_unregistered",
	Help: "Service Proposal unregistered",
})

var proposalExpired = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "propmon_proposal_expired",
	Help: "Service Proposal expired",
})

var proposalInvalid = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "propmon_proposal_invalid",
	Help: "Service Proposal invalid",
})

var proposalCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_proposal_count",
	Help: "Service Proposal count",
}, []string{"service_type"})

var providerCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_provider_count",
	Help: "Provider count",
}, []string{"country", "node_type"})

var natsBytesReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "propmon_nats_bytes_rx",
	Help: "Number of bytes received by NATS listener",
}, []string{"subject"})

var quality = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_quality",
	Help: "Average quality score for country and node type",
}, []string{"country", "node_type"})

var latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_latency",
	Help: "Average latency for country and node type",
}, []string{"country", "node_type"})

var bandwidth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_bandwidth",
	Help: "Average bandwidth for country and node type",
}, []string{"country", "node_type"})

var uptime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "propmon_uptime",
	Help: "Average uptime for country and node type",
}, []string{"country", "node_type"})

func init() {
	Registry.MustRegister(
		proposalRegistered,
		proposalPing,
		proposalUnregistered,
		proposalExpired,
		proposalInvalid,
		proposalCount,
		providerCount,
		natsBytesReceived,
		quality,
		latency,
		bandwidth,
		uptime,
	)
}

func ProposalPing() {
	proposalPing.Inc()
}

func ProposalRegistered() {
	proposalRegistered.Inc()
}

func ProposalInvalid() {
	proposalInvalid.Inc()
}

func ProposalUnregistered() {
	proposalUnregistered.Inc()
}

func ActiveProposals(proposals []*proposal.Proposal) {
	serviceTypes := make(map[string]int)

	for _, p := range proposals {
		serviceTypes[p.ServiceType]++
	}

	for serviceType, n := range serviceTypes {
		proposalCount.WithLabelValues(serviceType).Set(float64(n))
	}
}

func ActiveProviders(providers []*proposal.Provider) {
	totals := make(map[providerLabel]int)
	qualities := make(map[providerLabel]*proposal.Quality)

	for _, p := range providers {
		label := providerLabel{
			Country:  p.Location.Country,
			NodeType: p.Location.IpType,
		}
		totals[label]++

		if p.Quality == nil {
			continue
		}

		if qualities[label] == nil {
			qualities[label] = &proposal.Quality{}
		}
		qualities[label].Quality += p.Quality.Quality
		qualities[label].Latency += p.Quality.Latency
		qualities[label].Bandwidth += p.Quality.Bandwidth
		qualities[label].Uptime += p.Quality.Uptime
	}

	for label, count := range totals {
		providerCount.WithLabelValues(label.Country, label.NodeType).Set(float64(count))

		if qualities[label].Quality != 0 {
			quality.WithLabelValues(label.Country, label.NodeType).Set(qualities[label].Quality / float64(count))
		}
		if qualities[label].Quality != 0 {
			latency.WithLabelValues(label.Country, label.NodeType).Set(qualities[label].Latency / float64(count))
		}
		if qualities[label].Bandwidth != 0 {
			bandwidth.WithLabelValues(label.Country, label.NodeType).Set(qualities[label].Bandwidth / float64(count))
		}
		if qualities[label].Uptime != 0 {
			uptime.WithLabelValues(label.Country, label.NodeType).Set(qualities[label].Uptime / float64(count))
		}
	}
}

func UpdateMetrics(repository *proposal.Repository, expired int) {
	ActiveProposals(repository.Proposals())
	ActiveProviders(repository.Providers())
	proposalExpired.Add(float64(expired))
}

func NatsMsgReceived(msg *nats.Msg) {
	natsBytesReceived.WithLabelValues(msg.Subject).Add(float64(len(msg.Data)))
}
