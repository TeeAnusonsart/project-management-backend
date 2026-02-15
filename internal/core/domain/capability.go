package domain

type Capability struct {
	ID             uint
	CapabilityType string
	ControlType	string
	Widgets        []Widget
}

type CapabilityRepository interface {
	FindByType(capType string) (*Capability, error)
	FindByTypeAndControl(capType string, controlType string) (*Capability, error)
}
