package expiration

import (
	"sync"
	"time"

	"github.com/sch8ill/propmon/metrics"
	"github.com/sch8ill/propmon/proposal"
)

type Service struct {
	repository *proposal.Repository
	jobDelay   time.Duration
	stopCh     chan struct{}
	waitGroup  sync.WaitGroup
}

func NewExpirationService(repository *proposal.Repository, jobDelay time.Duration) *Service {
	return &Service{
		repository: repository,
		jobDelay:   jobDelay,
		stopCh:     make(chan struct{}),
	}
}

func (e *Service) Start() {
	e.waitGroup.Add(1)
	go e.run()
}

func (e *Service) Stop() {
	close(e.stopCh)
	e.waitGroup.Wait()
}

func (e *Service) run() {
	defer e.waitGroup.Done()

	for {
		select {
		case <-e.stopCh:
			return

		default:
			if e.repository.CountProposals() > 0 {
				expired := e.repository.RemoveExpired()
				metrics.ReportStatus(e.repository, expired)
			}
			time.Sleep(e.jobDelay)
		}
	}
}
