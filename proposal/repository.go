package proposal

import (
	"sync"
	"time"
)

type Repository struct {
	proposalLifetime time.Duration
	proposals        map[string]proposalRecord
	mu               sync.RWMutex
}

type proposalRecord struct {
	proposal Proposal
	expires  time.Time
}

func NewProposalRepository(proposalLifetime time.Duration) *Repository {
	return &Repository{
		proposalLifetime: proposalLifetime,
		proposals:        make(map[string]proposalRecord),
	}
}

func (pr *Repository) Store(p Proposal) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.proposals[p.ServiceKey()] = proposalRecord{
		proposal: p,
		expires:  time.Now().Add(pr.proposalLifetime),
	}
}

func (pr *Repository) Get(key string) Proposal {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	return pr.proposals[key].proposal
}

func (pr *Repository) Exists(key string) bool {
	if _, ok := pr.proposals[key]; ok {
		return true
	}

	return false
}

func (pr *Repository) Remove(key string) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	delete(pr.proposals, key)
}

func (pr *Repository) Proposals() []Proposal {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var proposals []Proposal

	for _, rcd := range pr.proposals {
		proposals = append(proposals, rcd.proposal)
	}

	return proposals
}

func (pr *Repository) Providers() []Provider {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var providers = make(map[string]Provider)

	for _, rcd := range pr.proposals {
		p := rcd.proposal

		if _, ok := providers[p.ProviderID]; ok {
			provider := providers[p.ProviderID]
			provider.Services = append(provider.Services, Service{
				ServiceType:    p.ServiceType,
				Compatibility:  p.Compatibility,
				Contacts:       p.Contacts,
				AccessPolicies: p.AccessPolicies,
			})
			providers[p.ProviderID] = provider
			continue
		}

		providers[p.ProviderID] = Provider{
			ID:       p.ProviderID,
			Location: p.Location,
			Quality:  p.Quality,
			Services: []Service{
				{
					ServiceType:    p.ServiceType,
					Compatibility:  p.Compatibility,
					Contacts:       p.Contacts,
					AccessPolicies: p.AccessPolicies,
				},
			},
		}
	}

	var providerSlice []Provider
	for _, p := range providers {
		providerSlice = append(providerSlice, p)
	}

	return providerSlice
}

func (pr *Repository) CountProposals() int {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	return len(pr.proposals)
}

func (pr *Repository) CountProviders() int {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	providers := make(map[string]bool)

	for _, rcd := range pr.proposals {
		providers[rcd.proposal.ProviderID] = true
	}

	return len(providers)
}

func (pr *Repository) RemoveExpired() int {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	var expired int

	for key, record := range pr.proposals {
		if time.Now().After(record.expires) {
			delete(pr.proposals, key)
			expired++
		}
	}

	return expired
}
