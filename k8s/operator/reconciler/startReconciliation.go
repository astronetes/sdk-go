package reconciler

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// setStatusToUnknown is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) startReconciliation(ctx context.Context, c client.Client, cfg Config, req ctrl.Request) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	var obj S

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	// Let's just set the status as Unknown when no status are available
	if obj.AstronetesStatus().Conditions == nil || len(obj.AstronetesStatus().Conditions) == 0 {
		meta.SetStatusCondition(
			&obj.AstronetesStatus().Conditions,
			metav1.Condition{
				Type:    typeReadyResource,
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
