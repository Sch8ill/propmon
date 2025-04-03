package proposal

type Provider struct {
	ID       string
	Location Location
	Quality  *Quality
	Services []Service
}

type Service struct {
	ServiceType    string
	Compatibility  int
	Contacts       []Contact
	AccessPolicies []AccessPolicy
}
