package reconciler

import (
	"github.com/astronetes/sdk-go/k8s/operator/provider"
	"time"
)

type Config struct {
	Timeout          *time.Duration                                `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RequeueDelayTime []time.Duration                               `json:"requeueDelayTime,omitempty" yaml:"requeueDelayTime,omitempty"`
	Meta             map[string]any                                `json:"meta,omitempty" yaml:"meta,omitempty"`
	Providers        map[provider.GroupID]provider.ProvidersConfig `json:"providers,omitempty"  yaml:"providers,omitempty"`
}
