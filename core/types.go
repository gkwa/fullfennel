package core

type InstanceStatus struct {
	InstanceID string `json:"instanceId"`
	State      string `json:"state"`
}

type EC2StatusChecker interface {
	GetEC2Status() (InstanceStatus, error)
}
