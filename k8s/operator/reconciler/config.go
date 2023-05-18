package reconciler

import (
	"time"

	"github.com/astronetes/sdk-go/k8s/operator/provider"
)

const defaultRequeueTime = 5 * time.Second

type Config struct {
	Timeout          *time.Duration                                `mapstructure:"timeout,omitempty" `
	RequeueDelayTime []time.Duration                               `mapstructure:"requeueDelayTime,omitempty"`
	MaxConditions    int                                           `mapstructure:"maxConditions,omitempty"`
	Meta             map[string]any                                `mapstructure:"meta,omitempty"`
	Providers        map[provider.GroupID]provider.ProvidersConfig `mapstructure:"providers,omitempty"`
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
