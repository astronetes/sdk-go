package controller

import (
	"context"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *astronetes[S]) reconcileReady(ctx context.Context, c client.Client, cfg config.Phase, obj S) Result {
	log := log.FromContext(ctx)
	log.Info("Reconciliation for phase ready")
	span := trace.SpanFromContext(ctx)
	_, span = span.TracerProvider().Tracer(r.ID).Start(ctx, "reconcile-ready")
	defer span.End()
	span.AddEvent("The resource is ready")
	obj.AstronetesStatus().Conditions[0].Status = metav1.ConditionFalse
	obj.AstronetesStatus().SetReady(true)
	return Completed("Resource is created successfully")
}
