package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/sch8ill/propmon/proposal"
)

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

func init() {
	registry.MustRegister(proposalRegistered, proposalPing, proposalUnregistered, proposalExpired, proposalInvalid, proposalCount, providerCount)
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

func ActiveProposals(proposals []proposal.Proposal) {
	serviceTypes := make(map[string]int)

	for _, p := range proposals {
		serviceTypes[p.ServiceType]++
	}

	for serviceType, n := range serviceTypes {
		proposalCount.WithLabelValues(serviceType).Set(float64(n))
	}
}

func ActiveProviders(providers []proposal.Provider) {
	labels := make(map[providerLabel]int)

	for _, p := range providers {
		l := providerLabel{
			Country:  p.Location.Country,
			NodeType: p.Location.IpType,
		}
		labels[l]++
	}

	for label, count := range labels {
		providerCount.WithLabelValues(label.Country, label.NodeType).Set(float64(count))
	}
}

func ReportStatus(repository *proposal.Repository, expired int) {
	ActiveProposals(repository.Proposals())
	ActiveProviders(repository.Providers())
	proposalExpired.Add(float64(expired))

	log.Info().Int("proposals", repository.CountProposals()).Int("providers", repository.CountProviders()).Int("expired", expired).Msg("removed expired proposals")
}
