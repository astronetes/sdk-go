package reconciler

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// handleCreation is a function of type subreconciler.FnWithRequest.
func (r *reconciler[S]) handleCreation(ctx context.Context, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the latest version of the resource
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	log.Info("Performing Reconciliation operations for the resource")
	if meta.IsStatusConditionPresentAndEqual(obj.ReconcilableStatus().Conditions, ConditionTypeReady,
		metav1.ConditionUnknown) {
		obj.ReconcilableStatus().Attempts += 1
	} else {
		obj.ReconcilableStatus().Attempts = 1
	}

	condition := meta.FindStatusCondition(obj.ReconcilableStatus().Conditions, ConditionTypeReady)

	// Condition doesn't exist and must be created
	if condition == nil {
		obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
			Type:    ConditionTypeReady,
			Status:  metav1.ConditionTrue,
			Reason:  ConditionReasonReconciling,
			Message: MessageReconciliationInProcess,
		})
	}

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, ErrorUpdatingStatus)
		return RequeueWithError(err)
	}

	// Perform all operations required before remove the finalizer and allow
	// the Kubernetes API to remove the custom resource.
	res, err := r.subreconciler.HandleReconciliation(ctx, obj)
	if err != nil {
		// Set Ready condition status to False
		condition := meta.FindStatusCondition(obj.ReconcilableStatus().Conditions, ConditionTypeReady)
		condition.Status = metav1.ConditionFalse
	}
	if updateStatusErr := r.Status().Update(ctx, obj); updateStatusErr != nil {
		log.Error(updateStatusErr, ErrorUpdatingStatus)

		return RequeueWithError(updateStatusErr)
	}
	if err != nil {
		r.SetReconcilingMessage(ctx, obj, err.Error())
		log.Error(err, "Failed to perform creation operations")
	}
	if ShouldRequeue(res, err) {
		if res.RequeueAfter == 0 {
			res.RequeueAfter = r.config.GetRequeueTimeForAttempt(int(obj.ReconcilableStatus().Attempts))
		}
		return res, err
	}

	return ContinueReconciling()
}
