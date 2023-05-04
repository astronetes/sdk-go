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
	Timeout       *time.Duration        `yaml:"timeout,omitempty"`
	MaxConditions int32                 `yaml:"maxConditions,omitempty"`
	RequeueAfter  *time.Duration        `yaml:"requeueAfter,omitempty"`
	Phases        []ReconciliationPhase `yaml:"phases,omitempty"`
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
	Timeout      *time.Duration  `yaml:"timeout,omitempty"`
	Backoff      []time.Duration `yaml:"backoff,omitempty"`
	RequeueAfter *time.Duration  `yaml:"requeueAfter,omitempty"`
	Meta         map[string]any  `yaml:"meta,omitempty"`
}

type Monitoring struct {
	Address  string `yaml:"address,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Phase struct {
	Timeout       *time.Duration  `yaml:"timeout,omitempty"`
	MaxConditions int32           `yaml:"maxConditions,omitempty"`
	Backoff       []time.Duration `yaml:"backoff,omitempty"`
	RequeueAfter  *time.Duration  `yaml:"requeueAfter,omitempty"`
	Meta          map[string]any  `yaml:"meta,omitempty"`
}

func (ctrl Controller) GetConfigForReconciliationPhase(code v1.PhaseCode) Phase {
	out := Phase{
		Timeout:       ctrl.Timeout,
		MaxConditions: ctrl.MaxConditions,
		Backoff:       []time.Duration{},
		RequeueAfter:  ctrl.RequeueAfter,
		Meta:          map[string]any{},
	}
	phase := ctrl.getConfigForReconciliationPhase(code)
	if phase.Backoff != nil {
		out.Backoff = phase.Backoff
	}
	if phase.RequeueAfter != nil {
		out.RequeueAfter = phase.RequeueAfter
	}
	if phase.Meta != nil {
		out.Meta = phase.Meta
	}
	if phase.Timeout != nil {
		out.Timeout = phase.Timeout
	}
	return out
}

type ReconciliationPhaseSettings struct{}
