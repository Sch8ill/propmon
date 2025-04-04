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
	proposal *Proposal
	expires  time.Time
}

func NewProposalRepository(proposalLifetime time.Duration) *Repository {
	return &Repository{
		proposalLifetime: proposalLifetime,
		proposals:        make(map[string]proposalRecord),
	}
}

func (r *Repository) Store(p *Proposal) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.proposals[p.ServiceKey()] = proposalRecord{
		proposal: p,
		expires:  time.Now().Add(r.proposalLifetime),
	}
}

func (r *Repository) Get(key string) *Proposal {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.proposals[key].proposal
}

func (r *Repository) Exists(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.proposals[key]; ok {
		return true
	}

	return false
}

func (r *Repository) Remove(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.proposals, key)
}

func (r *Repository) Renew(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rcd := r.proposals[id]
	rcd.expires = time.Now().Add(r.proposalLifetime)
	r.proposals[id] = rcd
}

func (r *Repository) RenewOrStore(p *Proposal) {
	id := p.ServiceKey()

	if r.Exists(id) {
		r.Renew(id)
	} else {
		r.Store(p)
	}
}

func (r *Repository) Proposals() []*Proposal {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// TODO: create fixed size array once null pointer dereference in metrics.ActiveProposals is fixed
	var proposals []*Proposal

	for _, rcd := range r.proposals {
		proposals = append(proposals, rcd.proposal)
	}

	return proposals
}

func (r *Repository) Match(filter *Proposal, max int) []*Proposal {
	var matches []*Proposal

	for _, p := range r.Proposals() {
		if filter.ProviderID != "" && filter.ProviderID != p.ProviderID {
			continue
		}

		if filter.ServiceType != "" && filter.ServiceType != p.ServiceType {
			continue
		}

		if filter.Location.Country != "" && filter.Location.Country != p.Location.Country {
			continue
		}

		if filter.Location.IpType != "" && filter.Location.IpType != p.Location.IpType {
			continue
		}

		matches = append(matches, p)
		if len(matches) >= max {
			break
		}
	}

	return matches
}

func (r *Repository) Providers() []*Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	providers := make(map[string]*Provider)

	for _, rcd := range r.proposals {
		p := rcd.proposal

		if provider, ok := providers[p.ProviderID]; ok {
			provider.Services = append(provider.Services, Service{
				ServiceType:    p.ServiceType,
				Compatibility:  p.Compatibility,
				Contacts:       p.Contacts,
				AccessPolicies: p.AccessPolicies,
			})
			continue
		}

		providers[p.ProviderID] = &Provider{
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

	var providerSlice []*Provider
	for _, p := range providers {
		providerSlice = append(providerSlice, p)
	}

	return providerSlice
}

func (r *Repository) Countries() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	countries := make(map[string]struct{})

	for _, rcd := range r.proposals {
		countries[rcd.proposal.Location.Country] = struct{}{}
	}

	var countriesSlice []string
	for country := range countries {
		countriesSlice = append(countriesSlice, country)
	}

	return countriesSlice
}

func (r *Repository) CountProposals() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.proposals)
}

func (r *Repository) CountProviders() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	providers := make(map[string]bool)

	for _, rcd := range r.proposals {
		providers[rcd.proposal.ProviderID] = true
	}

	return len(providers)
}

func (r *Repository) UpdateQuality(qualityData map[string]*Quality) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, quality := range qualityData {
		if rcd, ok := r.proposals[id]; ok {
			rcd.proposal.Quality = quality
			r.proposals[id] = rcd
		}
	}
}

func (r *Repository) RemoveExpired() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	var expired int

	for key, record := range r.proposals {
		if time.Now().After(record.expires) {
			delete(r.proposals, key)
			expired++
		}
	}

	return expired
}
