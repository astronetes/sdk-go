package controller

import (
	"context"
	"fmt"

	v12 "github.com/astronetes/sdk-go/k8s/operator/api/v1"

	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *astronetes[S]) reconcileFailed(ctx context.Context, client rClient.Client, cfg config.Phase, obj S) Result {
	log := log.FromContext(ctx)
	log.Info("Reconciliation for phase failed")
	span := trace.SpanFromContext(ctx)
	_, span = span.TracerProvider().Tracer(r.ID).Start(ctx, "reconcile-failed")
	defer span.End()
	if obj.Status().State != v12.FailedPhase {
		obj.Status().Conditions[0].Status = v1.ConditionFalse
		span.AddEvent("The execution cannot continue due to errors'")
		return Completed("The execution cannot continue due to errors")
	}
	if obj.GetDeletionTimestamp().IsZero() {
		msg := fmt.Sprintf("Reesume the creation of the resources")
		span.AddEvent(msg)
		obj.Status().Next(r.Dispatcher.InitialCreationPhaseCode, ResumeCreationEvent, "Resume the creation of the resources")
	} else {
		msg := fmt.Sprintf("Reesume the deletion of the resources")
		span.AddEvent(msg)
		obj.Status().Next(r.Dispatcher.InitialDeletionPhaseCode, ResumeDeletionEvent, "Resume the deletion of the resources")
	}
	return OK("resume execution")
}
