package controller

import (
	"context"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *astronetes[S]) reconcileTerminating(ctx context.Context, c_ client.Client, cfg config.Phase, obj S) Result {
	log := log.FromContext(ctx)
	log.Info("Reconciliation for phase terminating")
	span := trace.SpanFromContext(ctx)
	ctx, span = span.TracerProvider().Tracer(r.ID).Start(ctx, "reconcile-terminating")
	defer span.End()
	span.AddEvent("IngressController/Terminating")
	if controllerutil.ContainsFinalizer(obj, r.FinalizerName) {
		obj.AstronetesStatus().Next(r.Dispatcher.InitialDeletionPhaseCode, FinalizerExists, "Ready to start the destruction of the resources")
	} else {
		obj.AstronetesStatus().Next(r.Dispatcher.InitialDeletionPhaseCode, MissingFinalizer, "Already terminated")
	}
	return OK("ingress controller was updated successfully")
}
