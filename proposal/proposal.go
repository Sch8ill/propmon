package proposal

type Proposal struct {
	Format         string         `json:"format"`
	Compatibility  int            `json:"compatibility"`
	ProviderID     string         `json:"provider_id"`
	ServiceType    string         `json:"service_type"`
	Location       Location       `json:"location"`
	Contacts       []Contact      `json:"contacts"`
	Quality        Quality        `json:"quality"`
	AccessPolicies []AccessPolicy `json:"access_policies,omitempty"`
}

func (p Proposal) ServiceKey() string {
	return p.ProviderID + "." + p.ServiceType
}

type Location struct {
	Continent string `json:"continent"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Asn       int    `json:"asn"`
	Isp       string `json:"isp"`
	IpType    string `json:"ip_type"`
}

type Contact struct {
	Type       string            `json:"type"`
	Definition ContactDefinition `json:"definition"`
}

type ContactDefinition struct {
	BrokerAddresses []string `json:"broker_addresses"`
}

type AccessPolicy struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type Quality struct {
	Quality   float64 `json:"quality"`
	Latency   float64 `json:"latency"`
	Bandwidth float64 `json:"bandwidth"`
	Uptime    float64 `json:"uptime"`
}
