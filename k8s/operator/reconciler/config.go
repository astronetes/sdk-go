package reconciler

import (
	"time"

	"github.com/astronetes/sdk-go/k8s/operator/provider"
)

const defaultRequeueTime = 5 * time.Second

type Config struct {
	Timeout          *time.Duration                                `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RequeueDelayTime []time.Duration                               `json:"requeueDelayTime,omitempty" yaml:"requeueDelayTime,omitempty"`
	MaxConditions    int                                           `json:"maxConditions,omitempty" yaml:"maxConditions,omitempty"`
	Meta             map[string]any                                `json:"meta,omitempty" yaml:"meta,omitempty"`
	Providers        map[provider.GroupID]provider.ProvidersConfig `json:"providers,omitempty"  yaml:"providers,omitempty"`
}

func (c *Config) GetRequeueTimeForAttempt(attempt int) time.Duration {
	if len(c.RequeueDelayTime) == 0 {
		return defaultRequeueTime
	}
	if len(c.RequeueDelayTime) < attempt {
		return c.RequeueDelayTime[len(c.RequeueDelayTime)-1]
	}
	return c.RequeueDelayTime[attempt-1]
}
