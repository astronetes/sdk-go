package reconciler

var (
	// ConditionTypeReady represents the status of the Deployment reconciliation
	ConditionTypeReady = "Ready"
	// ConditionTypeDeleted represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	ConditionTypeDeleted = "Deleted"
)

var (
	ConditionReasonReconciling = "Reconciling"
	ConditionReasonReconciled  = "Reconciled"
	ConditionReasonDeleting    = "Deleting"
	ConditionReasonFinalizing  = "Finalizing"
)

var (
	MessageReconciliationInProcess = "Reconciliation in process"
	MessageReconciliationCompleted = "Reconciliation completed"
)

var (
	ErrorUpdatingStatus = "Failed to update resource status"
)
