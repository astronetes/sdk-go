package reconciler

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// handleCreation is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) handleCreation(ctx context.Context, c client.Client, cfg Config, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	log.Info("Performing Reconciliation operations for tje resource")
	obj.ReconcilableStatus().SetStatusCondition(
		metav1.Condition{
			Type:    typeReadyResource,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: fmt.Sprintf("Performing reconciling operations for the custom resource: %s ", obj.GetName()),
		})

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update resource status")
		return RequeueWithError(err)
	}

	// Perform all operations required before remove the finalizer and allow
	// the Kubernetes API to remove the custom resource.
	if _, err := r.doCreationOperationForResource(ctx, c, cfg, obj); err != nil {
		log.Error(err, "Failed to perform creation operations")
		return RequeueWithError(err)
	}

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "Failed to re-fetch resource")
		return RequeueWithError(err)
	}

	obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
		Type:    typeReadyResource,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("Reconciling operations for custom resource %s name were successfully accomplished", obj.GetName()),
	})

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update resource status")
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
