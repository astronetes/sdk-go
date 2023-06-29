package reconciler

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// handleDeletion is a function of type reconciler.FnWithRequest.
func (r *reconciler[S]) handleDeletion(ctx context.Context, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the latest version of the resource
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	// Check if the resource instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isResourceMarkedToBeDeleted := obj.GetDeletionTimestamp() != nil
	if isResourceMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(obj, r.finalizerName) {
			log.Info("Performing Finalizer Operations for resource before delete CR")

			if meta.IsStatusConditionPresentAndEqual(obj.ReconcilableStatus().Conditions, ConditionTypeDeleted,
				metav1.ConditionUnknown) {
				obj.ReconcilableStatus().Attempts += 1
			} else {
				obj.ReconcilableStatus().Attempts = 1
			}

			// Let's add here an status "Downgrade" to define that this resource begin its process to be terminated.
			obj.ReconcilableStatus().SetStatusCondition(
				metav1.Condition{
					Type:    ConditionTypeDeleted,
					Status:  metav1.ConditionUnknown,
					Reason:  ConditionReasonFinalizing,
					Message: fmt.Sprintf("Performing finalizer operations for the custom resource: %s ", obj.GetName()),
				})

			if err := r.Status().Update(ctx, obj); err != nil {
				log.Error(err, "Failed to update resource status")
				return RequeueWithError(err)
			}

			// Perform all operations required before remove the finalizer and allow
			// the Kubernetes API to remove the custom resource.
			res, err := r.subreconciler.HandleDeletion(ctx, obj)
			if updateStatusErr := r.Status().Update(ctx, obj); updateStatusErr != nil {
				log.Error(updateStatusErr, "Failed to update resource status")
				return RequeueWithError(updateStatusErr)
			}
			if err != nil {
				log.Error(err, "Failed to perform creation operations")
			}
			if ShouldRequeue(res, err) {
				if res.RequeueAfter == 0 {
					res.RequeueAfter = r.config.GetRequeueTimeForAttempt(int(obj.ReconcilableStatus().Attempts))
				}
				return res, err
			}

			// Re-fetch the resource Custom Resource before update the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raise the issue "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
				log.Error(err, "Failed to re-fetch resource")
				return RequeueWithError(err)
			}

			obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
				Type:   ConditionTypeDeleted,
				Status: metav1.ConditionTrue,
				Reason: ConditionReasonFinalizing,
				Message: fmt.Sprintf("Finalizer operations for custom resource %s name were successfully "+
					"accomplished", obj.GetName()),
			})

			if err := r.Status().Update(ctx, obj); err != nil {
				log.Error(err, "Failed to update the resource status")
				return RequeueWithError(err)
			}

			log.Info("Removing Finalizer for resource after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(obj, r.finalizerName); !ok {
				log.Error(nil, "Failed to remove finalizer for resource")
				return Requeue()
			}

			if err := r.Update(ctx, obj); err != nil {
				log.Error(err, "Failed to remove finalizer for resource")
				return RequeueWithError(err)
			}
		}
		return DoNotRequeue()
	}

	return ContinueReconciling()
}
