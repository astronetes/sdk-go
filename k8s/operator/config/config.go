package config

import (
	"time"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Config struct {
	Controllers map[string]Controller `yaml:"controllers,omitempty"`
	Monitoring  Monitoring            `yaml:"monitoring,omitempty"`
}

type Controller struct {
	Timeout             time.Duration         `yaml:"timeout,omitempty"`
	MaxConditions       int32                 `yaml:"max_conditions,omitempty"`
	DefaultRequeueAfter *time.Duration        `yaml:"defaultRequeueAfter,omitempty"`
	Phases              []ReconciliationPhase `yaml:"phases,omitempty"`
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
	Name         v1.PhaseCode    `yaml:"name,omitempty"`
	Backoff      []time.Duration `yaml:"backoff,omitempty"`
	RequeueAfter *time.Duration  `yaml:"requeueAfter,omitempty"`
}

type Monitoring struct {
	Address  string `yaml:"address,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Phase struct {
	Timeout       time.Duration   `yaml:"timeout,omitempty"`
	MaxConditions int32           `yaml:"max_conditions,omitempty"`
	Backoff       []time.Duration `yaml:"backoff,omitempty"`
	RequeueAfter  *time.Duration  `yaml:"requeueAfter,omitempty"`
}

func (ctrl Controller) GetConfigForReconciliationPhase(code v1.PhaseCode) Phase {
	out := Phase{
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
