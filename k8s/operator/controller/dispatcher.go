package controller

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type ReconciliationEvent string

type Dispatcher[Obj any] interface {
	ThrowEvent(ctx context.Context, obj Obj, evt ReconciliationEvent, message string) error
	GetReconciliationPhase(phase v1.PhaseCode) (PhaseReconcile[Obj], error)
	IsOnDeletionPhase(phase v1.PhaseCode) bool
}
