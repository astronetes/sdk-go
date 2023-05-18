package reconciler

import (
	"context"
	errors2 "github.com/astronetes/sdk-go/k8s/operator/errors"
	corev1 "k8s.io/api/core/v1"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	log2 "github.com/astronetes/sdk-go/log"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Requeue returns a controller result pairing specifying to
// requeue with no error message implied. This returns no error.
func Requeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, nil }

// RequeueWithError returns a controller result pairing specifying to
// requeue with an error message.
func RequeueWithError(e error) (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, e }

// DoNotRequeue returns a controller result pairing specifying not to requeue.
func DoNotRequeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: false}, nil }

// ContinueReconciling indicates that the reconciliation block should continue by
// returning a nil result and a nil error
func ContinueReconciling() (*reconcile.Result, error) { return nil, nil }

// ShouldHaltOrRequeue returns true if reconciler result is not nil
// or the err is not nil. In theory, the error evaluation
// is not needed because ShouldRequeue handles it, but
// it's included in case ShouldHaltOrRequeue is called directly.
func ShouldHaltOrRequeue(r *ctrl.Result, err error) bool {
	return (r != nil) || ShouldRequeue(r, err)
}

// Evaluate returns the actual reconcile struct and error. Wrap helpers in
// this when returning from within the top-level Reconciler.
func Evaluate(r *reconcile.Result, e error) (reconcile.Result, error) {
	return *r, e
}

// ShouldRequeue returns true if the reconciler result indicates
// a requeue is required, or the error is not nil.
func ShouldRequeue(r *ctrl.Result, err error) bool {
	// if we get a nil value for result, we need to
	// fill it with an empty value which would not trigger
	// a requeue.

	res := r
	if r.IsZero() {
		res = &ctrl.Result{}
	}
	return res.Requeue || (err != nil)
}

type Reconciler[S v1.Resource] interface {
	Reconcile(ctx context.Context, req ctrl.Request, obj S) (ctrl.Result, error)
}

type reconciler[S v1.Resource] struct {
	client.Client
	config        Config
	finalizerName string
	Recorder      record.EventRecorder
	Tracer        trace.Tracer
	Scheme        *runtime.Scheme
	subReconciler Handler[S]
}

func New[S v1.Resource](id string, mgr manager.Manager, finalizerName string,
	config Config, subReconciler Handler[S],
) Reconciler[S] {
	return &reconciler[S]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		finalizerName: finalizerName,
		subReconciler: subReconciler,
		Recorder:      mgr.GetEventRecorderFor(id),
		Tracer:        otel.Tracer(id),
		config:        config,
	}
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

func (r *reconciler[S]) Reconcile(ctx context.Context, req ctrl.Request, obj S) (ctrl.Result, error) {
	log := log2.FromContext(ctx)
	// The list of subreconcilers for resource
	subreconcilersForResource := []func(ctx context.Context, req ctrl.Request, obj S) (*ctrl.Result, error){
		r.startReconciliation,
		r.addFinalizer,
		r.handleDeletion,
		r.handleCreation,
		r.updateStatus,
	}

	// Run all subreconcilers sequentially
	for _, f := range subreconcilersForResource {
		if res, err := f(ctx, req, obj); ShouldHaltOrRequeue(res, err) {
			if err != nil {
				switch x := err.(type) {
				case *errors2.ResourceError:
					r.Recorder.Eventf(obj, corev1.EventTypeWarning, string(x.Code()), "'%s', check documentation at '%s", x.Msg(), x.DocRef())
				case *errors2.ControllerError:
					r.Recorder.Eventf(obj, corev1.EventTypeWarning, string(x.Code()), "'%s', check documentation at '%s", x.Msg(), x.DocRef())
				default:
					r.Recorder.Event(obj, corev1.EventTypeWarning, "Error", err.Error())
				}
			}
			return Evaluate(res, err)
		}
		if err := r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "Failed to update the resource status")
			return Evaluate(RequeueWithError(err))
		}
	}

	return Evaluate(DoNotRequeue())
}

func (r *reconciler[S]) newObject() S {
	var objValue S
	return objValue
}
