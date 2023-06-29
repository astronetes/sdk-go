package reconciler

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// setStatusToUnknown is a function of type subreconciler.FnWithRequest.
func (r *reconciler[S]) startReconciliation(ctx context.Context, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	// Let's just set the status as Unknown when no status are available
	if obj.ReconcilableStatus().Conditions == nil || len(obj.ReconcilableStatus().Conditions) == 0 {
		obj.ReconcilableStatus().SetStatusCondition(
			metav1.Condition{
				Type:    ConditionTypeReady,
				Status:  metav1.ConditionUnknown,
				Reason:  "Reconciling",
				Message: "Starting reconciliation",
			},
		)
		if err := r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "Failed to update resource status")
			return RequeueWithError(err)
		}
	}

	return ContinueReconciling()
}
