package reconciler

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Handler[S v1.Resource] interface {
	Reconcile(ctx context.Context, obj S) (*ctrl.Result, error)
	Delete(ctx context.Context, obj S) (*ctrl.Result, error)
	SetDeletingMessage(ctx context.Context, obj S, msg string) error
	SetReconcilingMessage(ctx context.Context, obj S, msg string) error
}

type SubReconcilerHandler[S v1.Resource] struct {
	client.Client
	Manager  manager.Manager
	Config   Config
	Recorder record.EventRecorder
	Tracer   trace.Tracer
	Scheme   *runtime.Scheme
}

func (h *SubReconcilerHandler[S]) SetState(ctx context.Context, obj S, state v1.PhaseCode) error {
	log := log.FromContext(ctx)
	obj.ReconcilableStatus().State = state
	if err := h.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update status")
		return err
	}
	h.RecordEvent(obj, string(state), "Set status to '%s'", string(state))
	return nil
}

func (h *SubReconcilerHandler[S]) RecordEvent(obj S, reason string, msg string, args ...interface{}) {
	h.Recorder.Eventf(obj, corev1.EventTypeWarning, reason, msg, args...)
}

func (h *SubReconcilerHandler[S]) SetConditionMessageByType(ctx context.Context, obj S, conditionType, msg string) error {
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

	if err := h.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update object status")
		return err
	}
	h.RecordEvent(obj, string(msg), "Set message to '%s'", string(msg))
	return nil
}

func (h *SubReconcilerHandler[S]) SetDeletingMessage(ctx context.Context, obj S, msg string) error {
	return h.SetConditionMessageByType(ctx, obj, ConditionReasonDeleting, msg)
}

func (h *SubReconcilerHandler[S]) SetReconcilingMessage(ctx context.Context, obj S, msg string) error {
	return h.SetConditionMessageByType(ctx, obj, ConditionTypeReady, msg)
}
