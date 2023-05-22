package reconciler

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Handler[S v1.Resource] interface {
	Reconcile(ctx context.Context, obj S) (*ctrl.Result, error)
	Delete(ctx context.Context, obj S) (*ctrl.Result, error)
	// SetState(ctx context.Context, obj S, state v1.PhaseCode) error
	// RecordEvent(obj S, reason string, msg string, args ...interface{})
	// SetConditionMessageByType(ctx context.Context, obj S, conditionType, msg string) error
	// SetDeletingMessage(ctx context.Context, obj S, msg string) error
	// SetReconcilingMessage(ctx context.Context, obj S, msg string) error
}
