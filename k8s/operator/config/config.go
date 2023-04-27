package config

import (
	"time"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Config struct {
	Controllers map[string]Controller `json:"controllers,omitempty"`
	Monitoring  Monitoring            `json:"monitoring,omitempty"`
}

type Controller struct {
	Timeout             time.Duration         `json:"timeout,omitempty"`
	MaxConditions       int32                 `json:"max_conditions,omitempty"`
	DefaultRequeueAfter *time.Duration        `json:"defaultRequeueAfter,omitempty"`
	Phases              []ReconciliationPhase `json:"phases,omitempty"`
}

func (ctrl Controller) getConfigForReconciliationPhase(code v1.PhaseCode) ReconciliationPhase {
	for _, p := range ctrl.Phases {
		if p.Name == code {
			return p
		}
	}
	return ReconciliationPhase{}
}

type ReconciliationPhase struct {
	Name         v1.PhaseCode    `json:"name,omitempty"`
	Backoff      []time.Duration `json:"backoff,omitempty"`
	RequeueAfter *time.Duration  `json:"requeueAfter,omitempty"`
}

type Monitoring struct {
	Address  string `json:"address,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type PhaseSettings struct {
	Timeout       time.Duration   `json:"timeout,omitempty"`
	MaxConditions int32           `json:"max_conditions,omitempty"`
	Backoff       []time.Duration `json:"backoff,omitempty"`
	RequeueAfter  *time.Duration  `json:"requeueAfter,omitempty"`
}

func (ctrl Controller) GetConfigForReconciliationPhase(code v1.PhaseCode) PhaseSettings {
	out := PhaseSettings{
		Timeout:       ctrl.Timeout,
		MaxConditions: ctrl.MaxConditions,
		RequeueAfter:  ctrl.DefaultRequeueAfter,
	}
	phase := ctrl.getConfigForReconciliationPhase(code)
	if phase.Backoff != nil {
		out.Backoff = phase.Backoff
	}
	if phase.RequeueAfter != nil {
		out.RequeueAfter = phase.RequeueAfter
	}
	return out
}

type ReconciliationPhaseSettings struct{}
