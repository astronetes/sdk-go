package reconciler

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Subreconciler[S v1.Resource] struct {
	client.Client
	Manager  manager.Manager
	Config   Config
	Recorder record.EventRecorder
	Tracer   trace.Tracer
	Scheme   *runtime.Scheme
}

func (r *Subreconciler[S]) SetConditionMessageByType(ctx context.Context, obj S, conditionType, msg string) error {
	return setConditionMessageByType[S](ctx, r.Client, r.Recorder, obj, conditionType, msg)
}

func (r *Subreconciler[S]) SetDeletingMessage(ctx context.Context, obj S, msg string) error {
	return r.SetConditionMessageByType(ctx, obj, ConditionReasonDeleting, msg)
}

func (r *Subreconciler[S]) SetReconcilingMessage(ctx context.Context, obj S, msg string) error {
	return r.SetConditionMessageByType(ctx, obj, ConditionTypeReady, msg)
}

func (r *Subreconciler[S]) SetState(ctx context.Context, obj S, state v1.PhaseCode) error {
	log := log.FromContext(ctx)
	obj.ReconcilableStatus().State = state
	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update status")
		return err
	}
	r.RecordEvent(obj, string(state), "Set status to '%s'", string(state))
	return nil
}

func (r *Subreconciler[S]) RecordEvent(obj S, reason string, msg string, args ...interface{}) {
	r.Recorder.Eventf(obj, corev1.EventTypeWarning, reason, msg, args...)
}
