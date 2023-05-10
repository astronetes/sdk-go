package controller

import (
	"fmt"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Dispatcher struct {
	GetPhase func(code v1.PhaseCode) (PhaseReconcile[v1.Resource], error) {
	IsOnDeletionPhase func(code v1.PhaseCode) bool
	InitialPhaseCode  v1.PhaseCode
}

