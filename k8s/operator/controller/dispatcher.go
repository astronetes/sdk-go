package controller

import (
	"context"
	"fmt"

	"github.com/astronetes/sdk-go/k8s/operator/config"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
)

type Dispatcher[S v1.Resource] struct {
	Phase             func(code v1.PhaseCode) (PhaseReconcile[S], bool)
	IsOnDeletionPhase func(code v1.PhaseCode) bool
	InitialPhaseCode  v1.PhaseCode
}

func (m Dispatcher[S]) ReconcilePhase(ctx context.Context, code v1.PhaseCode, c client.Client, cfg config.Phase, obj S) (Result, error) {
	p, ok := m.Phase(code)
	if !ok {
		return Result{}, fmt.Errorf("unknown phase")
	}
	return p(ctx, c, cfg, obj), nil
}
