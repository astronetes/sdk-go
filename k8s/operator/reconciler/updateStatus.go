package reconciler

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// updateStatus is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) updateStatus(ctx context.Context, req ctrl.Request) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	var obj S

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}
	// The following implementation will update the status
	meta.SetStatusCondition(&obj.Status().Conditions, metav1.Condition{
		Type:    typeReadyResource,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("Creations of resources  for custom resource (%s) with was completed successfully", obj.GetName())})

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update Memcached status")
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
