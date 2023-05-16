package reconciler

import (
	"context"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Definitions to manage status conditions
const (
	// typeReadyResource represents the status of the Deployment reconciliation
	typeReadyResource = "Ready"
	// typeDegradedResource represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	typeDegradedResource = "Degraded"
)

type Reconciler[S v1.Resource] interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
}

type reconciler[S v1.Resource] struct {
	client.Client
	finalizerName                   string
	Recorder                        record.EventRecorder
	Tracer                          trace.Tracer
	Scheme                          *runtime.Scheme
	doDeletionOperationsForResource func(ctx context.Context, req ctrl.Request, obj S) error
	doCreationOperationForResource  func(ctx context.Context, req ctrl.Request, obj S) error
}

func New[S v1.Resource](id string, mgr manager.Manager, finalizerName string,
	config config.Controller) Reconciler[S] {
	return &reconciler[S]{
		Client:                          mgr.GetClient(),
		Scheme:                          mgr.GetScheme(),
		finalizerName:                   finalizerName,
		doDeletionOperationsForResource: nil,
		doCreationOperationForResource:  nil,
		Recorder:                        mgr.GetEventRecorderFor(id),
		Tracer:                          otel.Tracer(id),
	}
}

// setStatusToUnknown is a function of type subreconciler.FnWithRequest
func (r *reconciler[S]) startReconciliation(ctx context.Context, req ctrl.Request) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	var obj S

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}

	// Let's just set the status as Unknown when no status are available
	if obj.Status().Conditions == nil || len(obj.Status().Conditions) == 0 {
		meta.SetStatusCondition(
			&obj.Status().Conditions,
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

func (r *reconciler[S]) getLatest(ctx context.Context, req ctrl.Request, obj S) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("the resource could not be found. Ignoring since object must be deleted")
			return DoNotRequeue()
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get the resource")
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
