package operator

import (
	"time"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Settings struct {
	Timeout       time.Duration
	MaxConditions int
	PhaseOpts     map[v1.PhaseCode]interface{}
}
