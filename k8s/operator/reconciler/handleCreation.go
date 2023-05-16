package reconciler

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// handleCreation is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) handleCreation(ctx context.Context, req ctrl.Request) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	var obj S

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	// The following implementation will update the status
	meta.SetStatusCondition(
		&obj.Status().Conditions,
		metav1.Condition{
			Type:    typeReadyResource,
			Status:  metav1.ConditionFalse,
			Reason:  "Reconciled",
			Message: fmt.Sprintf("The custom resource (%s) has been created usccessfully", obj.GetName()),
		},
	)

	// Perform all operations required before remove the finalizer and allow
	// the Kubernetes API to remove the custom resource.
	if err := r.doCreationOperationForResource(ctx, req, obj); err != nil {
		log.Error(err, "Failed to perform finalizer operations")
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
