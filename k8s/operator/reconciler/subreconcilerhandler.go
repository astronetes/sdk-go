package reconciler

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	log := log.FromContext(ctx)

	condition := meta.FindStatusCondition(obj.ReconcilableStatus().Conditions, conditionType)

	// Condition doesn't exist and must be created
	if condition == nil {
		obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
			Type:    conditionType,
			Status:  metav1.ConditionTrue,
			Reason:  ConditionReasonReconciling,
			Message: msg,
		})

		// Condition exists and must be updated
	} else {
		condition.Message = msg
		meta.SetStatusCondition(
			&obj.ReconcilableStatus().Conditions,
			*condition,
		)
	}

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update object status")
		return err
	}
	r.RecordEvent(obj, string(msg), "Set message to '%s'", string(msg))
	return nil
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
