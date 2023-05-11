package controller

import (
	"context"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *astronetes[S]) reconcileDeleted(ctx context.Context, c client.Client, cfg config.Phase, obj S) Result {
	log := log.FromContext(ctx)
	log.Info("Reconciliation for phase deleted")
	span := trace.SpanFromContext(ctx)
	_, span = span.TracerProvider().Tracer(r.ID).Start(ctx, "reconcile-deleted")
	defer span.End()
	span.AddEvent("IngressController/Reconcile/Deleted")
	if controllerutil.ContainsFinalizer(obj, r.FinalizerName) {
		controllerutil.RemoveFinalizer(obj, r.FinalizerName)
	}
	if err := c.Update(ctx, obj); err != nil {
		log.Info(err.Error())
	}
	return Completed("ingress controller is deleted")
}
