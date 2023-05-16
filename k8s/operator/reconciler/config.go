package reconciler

import "time"

type Config struct {
	Timeout *time.Duration `yaml:"timeout,omitempty"`
}
