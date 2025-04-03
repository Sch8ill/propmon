package quality

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sch8ill/propmon/proposal"
)

type Service struct {
	oracle           *Oracle
	repository       *proposal.Repository
	proposalLifetime time.Duration
	interval         time.Duration
	stopCh           chan struct{}
	waitGroup        sync.WaitGroup
}

func NewQualityService(oracle *Oracle, repository *proposal.Repository, interval time.Duration, proposalLifetime time.Duration) *Service {
	return &Service{
		oracle:           oracle,
		repository:       repository,
		proposalLifetime: proposalLifetime,
		interval:         interval,
	}
}

func (s *Service) Start() {
	log.Debug().Msg("Starting quality service")
	s.waitGroup.Add(1)
	go s.run()
}

func (s *Service) Stop() {
	close(s.stopCh)
	s.waitGroup.Wait()
}

func (s *Service) run() {
	// wait until most proposals have been captured to not waste any quality entries
	time.Sleep(s.proposalLifetime)
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.stopCh:
			return

		default:
			if err := s.update(); err != nil {
				log.Warn().Err(err).Msg("Failed to update quality data")
			}
			time.Sleep(s.interval)
		}
	}
}

func (s *Service) update() error {
	qualityData, err := s.oracle.Quality()
	if err != nil {
		return fmt.Errorf("failed to fetch quality data: %w", err)
	}
	log.Debug().Msgf("Fetched %d quality entries", len(qualityData))
	s.repository.UpdateQuality(qualityData)

	return nil
}
