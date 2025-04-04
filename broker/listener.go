package broker

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/sch8ill/propmon/metrics"
	"github.com/sch8ill/propmon/proposal"
)

const (
	pingSubject       = "*.proposal-ping.v3"
	registerSubject   = "*.proposal-register.v3"
	unregisterSubject = "*.proposal-unregister.v3"
)

type Listener struct {
	brokerUrl  string
	conn       *nats.Conn
	repository *proposal.Repository
}

type Msg struct {
	Proposal *proposal.Proposal
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

	log.Info().Str("addr", l.brokerUrl).Msg("Connected to broker")

	if _, err := l.conn.Subscribe(pingSubject, l.onPing); err != nil {
		return err
	}

	if _, err := l.conn.Subscribe(registerSubject, l.onRegistration); err != nil {
		return err
	}

	if _, err := l.conn.Subscribe(unregisterSubject, l.onUnregistration); err != nil {
		return err
	}

	return nil
}

func (l *Listener) onPing(msg *nats.Msg) {
	metrics.NatsMsgReceived(msg)
	p, err := parseProposal(msg)
	if err != nil {
		metrics.ProposalInvalid()
		return
	}
	l.repository.RenewOrStore(p)
	metrics.ProposalPing()
}

func (l *Listener) onRegistration(msg *nats.Msg) {
	metrics.NatsMsgReceived(msg)
	p, err := parseProposal(msg)
	if err != nil {
		metrics.ProposalInvalid()
		return
	}
	l.repository.Store(p)
	metrics.ProposalRegistered()
}

func (l *Listener) onUnregistration(msg *nats.Msg) {
	metrics.NatsMsgReceived(msg)
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

func parseProposal(msg *nats.Msg) (*proposal.Proposal, error) {
	var brokerMsg Msg
	if err := json.Unmarshal(msg.Data, &brokerMsg); err != nil {
		return nil, fmt.Errorf("failed to parse proposal: %w", err)
	}

	return brokerMsg.Proposal, nil
}
