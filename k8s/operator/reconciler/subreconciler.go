package reconciler

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Definitions to manage status conditions
const (
	// typeReadyResource represents the status of the Deployment reconciliation
	typeReadyResource = "Ready"
	// typeDegradedResource represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	typeDegradedResource = "Degraded"
)

type Handler[S v1.Resource] interface {
	Reconcile(ctx context.Context, obj S) (*ctrl.Result, error)
	Delete(ctx context.Context, obj S) (*ctrl.Result, error)
}

type SubReconcilerHandler[S v1.Resource] struct {
	client.Client
	Config   Config
	Recorder record.EventRecorder
	Tracer   trace.Tracer
	Scheme   *runtime.Scheme
}

func (h *SubReconcilerHandler[S]) SetState(ctx context.Context, obj S, state v1.PhaseCode) error {
	log := log.FromContext(ctx)
	obj.ReconcilableStatus().State = state
	if err := h.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update Memcached status")
		return err
	}
	return nil
}
