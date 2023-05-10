package controller

import (
	"fmt"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Dispatcher struct {
	Phases            map[v1.PhaseCode]PhaseReconcile[v1.Resource]
	IsOnDeletionPhase func(code v1.PhaseCode) bool
	InitialPhaseCode  v1.PhaseCode
}

func (m Dispatcher) GetPhase(code v1.PhaseCode) (PhaseReconcile[v1.Resource], error) {
	p, ok := m.Phases[code]
	if !ok {
		return nil, fmt.Errorf("unknown phase")
	}
	return p, nil
}