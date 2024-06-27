package core

import (
	"fmt"
	"math/rand"
	"sync"
)

type RealEC2StatusChecker struct {
	instanceID string
}

func (c *RealEC2StatusChecker) GetEC2Status() (InstanceStatus, error) {
	return InstanceStatus{
		InstanceID: c.instanceID,
		State:      "running",
	}, nil
}

type MockEC2StatusChecker struct {
	instanceID   string
	states       []string
	currentIndex int
	rng          *rand.Rand
	mutex        sync.Mutex
}

func NewMockEC2StatusChecker(instanceID string, rng *rand.Rand) *MockEC2StatusChecker {
	states := []string{"running", "stopping", "stopped", "starting"}
	return &MockEC2StatusChecker{
		instanceID:   instanceID,
		states:       states,
		currentIndex: rng.Intn(len(states)),
		rng:          rng,
	}
}

func (m *MockEC2StatusChecker) GetEC2Status() (InstanceStatus, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.states) == 0 {
		return InstanceStatus{}, fmt.Errorf("no states defined")
	}

	state := m.states[m.currentIndex]
	m.currentIndex = (m.currentIndex + 1) % len(m.states)

	// time.Sleep(time.Duration(m.rng.Intn(1000)) * time.Millisecond)

	return InstanceStatus{
		InstanceID: m.instanceID,
		State:      state,
	}, nil
}
