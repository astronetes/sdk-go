package config

import (
	"github.com/astronetes/sdk-go/k8s/operator/reconciler"
)

type Config struct {
	Controllers map[string]reconciler.Config `mapstructure:"controllers,omitempty"`
	Monitoring  Monitoring                   `mapstructure:"monitoring,omitempty"`
	Namespace   string                       `mapstructure:"namespace,omitempty"`
}

type Monitoring struct {
	Address  string `mapstructure:"address,omitempty"`
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
}
