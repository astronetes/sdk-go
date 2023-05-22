package reconciler

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Requeue returns a controller result pairing specifying to
// requeue with no error message implied. This returns no error.
func Requeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, nil }

// RequeueWithError returns a controller result pairing specifying to
// requeue with an error message.
func RequeueWithError(e error) (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, e }

// DoNotRequeue returns a controller result pairing specifying not to requeue.
func DoNotRequeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: false}, nil }

// ContinueReconciling indicates that the reconciliation block should continue by
// returning a nil result and a nil error.
func ContinueReconciling() (*reconcile.Result, error) { return nil, nil }

// ShouldHaltOrRequeue returns true if reconciler result is not nil
// or the err is not nil. In theory, the error evaluation
// is not needed because ShouldRequeue handles it, but
// it's included in case ShouldHaltOrRequeue is called directly.
func ShouldHaltOrRequeue(r *ctrl.Result, err error) bool {
	return (r != nil) || ShouldRequeue(r, err)
}

// Evaluate returns the actual reconcile struct and error. Wrap helpers in
// this when returning from within the top-level Reconciler.
func Evaluate(r *reconcile.Result, e error) (reconcile.Result, error) {
	return *r, e
}

// ShouldRequeue returns true if the reconciler result indicates
// a requeue is required, or the error is not nil.
func ShouldRequeue(r *ctrl.Result, err error) bool {
	// if we get a nil value for result, we need to
	// fill it with an empty value which would not trigger
	// a requeue.

	res := r
	if r.IsZero() {
		res = &ctrl.Result{}
	}
	return res.Requeue || (err != nil)
}

func setConditionMessageByType[S v1.Resource](ctx context.Context, client client.Client, recorder record.EventRecorder,
	obj S, conditionType string, msg string,
) error {
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

	if err := client.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update object status")
		return err
	}
	recorder.Eventf(obj, corev1.EventTypeNormal, "Set message to '%s'", string(msg))
	return nil
}
