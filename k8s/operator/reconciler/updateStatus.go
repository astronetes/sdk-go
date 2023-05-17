package reconciler

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// updateStatus is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) updateStatus(ctx context.Context, c client.Client, cfg Config, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}
	// The following implementation will update the status
	obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
		Type:    typeReadyResource,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("Creations of resources  for custom resource (%s) with was completed successfully", obj.GetName()),
	})

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update Memcached status")
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
