package config

import (
	"github.com/astronetes/sdk-go/k8s/operator/reconciler"
)

type Config struct {
	Controllers map[string]reconciler.Config `yaml:"controllers,omitempty"`
	Monitoring  Monitoring                   `yaml:"monitoring,omitempty"`
}

type Monitoring struct {
	Address  string `yaml:"address,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
