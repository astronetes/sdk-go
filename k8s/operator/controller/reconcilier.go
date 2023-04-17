package controller

import (
	"context"

	"github.com/astronetes/sdk-go/log"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"helm.sh/helm/v3/pkg/time"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconciler struct {
	d *reconcilerDispatcher
}

type reconcilerDispatcher struct{}

func (d *reconcilerDispatcher) dispatch() ReconcilePhase {
	return nil
}

type reconcilerExecutor struct {
	Tracer trace.Tracer
}

func (re *reconcilerExecutor) Execute(ctx context.Context, phase ReconcilePhase, ctrl AstronetesController) {
	ctx, span := re.Tracer.Start(ctx, "ingress-controller-reconcile-error")
	defer span.End()
	l := log.FromContext(ctx)
	l.V(log.Info).Info("executing phase", "phase", phase.Info().name, "controller", ctrl.Name)
	span.AddEvent("",
		trace.WithAttributes(
			attribute.String("phase", phase.Info().name),
			attribute.String("controller", ctrl.Name()),
		),
		trace.WithTimestamp(time.Now().Time))

	res, err := phase.Reconcile(ctx, l)
	if err != nil {
		span.RecordError(err,
			trace.WithAttributes(
				attribute.String("phase", phase.Info().name),
				attribute.String("controller", ctrl.Name()),
			),
			trace.WithStackTrace(true),
		)
	}
	span.End()
}

type ReconcilerPhaseInfo struct {
	name        string
	description string
}

type ReconcilePhase interface {
	Info() ReconcilerPhaseInfo
	Reconcile(ctx context.Context, log logr.Logger) (reconcile.Result, error)
}
