package domain

type CapabilityRepository interface {
	FindByType(capType string) (*Capability, error)
}
