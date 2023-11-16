package broker

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/sch8ill/propmon/metrics"
	"github.com/sch8ill/propmon/proposal"
)

type Listener struct {
	brokerUrl  string
	conn       *nats.Conn
	repository *proposal.Repository
}

type Msg struct {
	Proposal proposal.Proposal
}

func NewListener(brokerUrl string, repository *proposal.Repository) *Listener {
	return &Listener{
		brokerUrl:  brokerUrl,
		repository: repository,
	}
}

func (l *Listener) Listen() error {
	conn, err := nats.Connect(l.brokerUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to broker: %w", err)
	}
	l.conn = conn

	log.Info().Str("addr", l.brokerUrl).Msg("connected to broker")

	if _, err := l.conn.Subscribe("*.proposal-ping.v3", l.onPing); err != nil {
		return err
	}

	if _, err := l.conn.Subscribe("*.proposal-register.v3", l.onRegistration); err != nil {
		return err
	}

	if _, err := l.conn.Subscribe("*.proposal-unregister.v3", l.onUnregistration); err != nil {
		return err
	}

	return nil
}

func (l *Listener) onPing(msg *nats.Msg) {
	p, err := parseProposal(msg)
	if err != nil {
		metrics.ProposalInvalid()
		return
	}
	l.repository.Store(p)
	metrics.ProposalPing()
}

func (l *Listener) onRegistration(msg *nats.Msg) {
	p, err := parseProposal(msg)
	if err != nil {
		metrics.ProposalInvalid()
		return
	}
	l.repository.Store(p)
	metrics.ProposalRegistered()
}

func (l *Listener) onUnregistration(msg *nats.Msg) {
	p, err := parseProposal(msg)
	if err != nil {
		metrics.ProposalInvalid()
		return
	}
	l.repository.Remove(p.ServiceKey())
	metrics.ProposalUnregistered()
}

func (l *Listener) Shutdown() {
	l.conn.Close()
}

func parseProposal(msg *nats.Msg) (proposal.Proposal, error) {
	brokerMsg := &Msg{}

	if err := json.Unmarshal(msg.Data, brokerMsg); err != nil {
		return proposal.Proposal{}, fmt.Errorf("failed to parse proposal: %w", err)
	}

	return brokerMsg.Proposal, nil
}
