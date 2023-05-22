package reconciler

import (
	"context"

	errors2 "github.com/astronetes/sdk-go/k8s/operator/errors"
	corev1 "k8s.io/api/core/v1"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	log2 "github.com/astronetes/sdk-go/log"

	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler[S v1.Resource] interface {
	Reconcile(ctx context.Context, req ctrl.Request, obj S) (ctrl.Result, error)
}

func NewReconciler[S v1.Resource](id string, mgr manager.Manager, finalizerName string,
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

type reconciler[S v1.Resource] struct {
	client.Client
	config        Config
	finalizerName string
	Recorder      record.EventRecorder
	Tracer        trace.Tracer
	Scheme        *runtime.Scheme
	subReconciler Handler[S]
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
	r.Recorder.Event(obj, corev1.EventTypeNormal, ConditionReasonReconciling, "Starting the reconciliation process.")
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
	r.Recorder.Event(obj, corev1.EventTypeNormal, ConditionReasonReconciling, MessageReconciliationCompleted)
	return Evaluate(DoNotRequeue())
}

func (r *reconciler[S]) newObject() S {
	var objValue S
	return objValue
}

func (r *reconciler[S]) RecordEvent(obj S, reason string, msg string, args ...interface{}) {
	r.Recorder.Eventf(obj, corev1.EventTypeWarning, reason, msg, args...)
}

func (r *reconciler[S]) SetConditionMessageByType(ctx context.Context, obj S, conditionType, msg string) error {
	log := log.FromContext(ctx)

	condition := meta.FindStatusCondition(obj.ReconcilableStatus().Conditions, conditionType)

	// Condition doesn't exist and must be created
	if condition == nil {
		obj.ReconcilableStatus().SetStatusCondition(metav1.Condition{
			Type:    conditionType,
			Status:  metav1.ConditionTrue,
			Reason:  ConditionReasonReconciling,
			Message: msg,
		})

		// Condition exists and must be updated
	} else {
		condition.Message = msg
		meta.SetStatusCondition(
			&obj.ReconcilableStatus().Conditions,
			*condition,
		)
	}

	if err := r.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update object status")
		return err
	}
	r.RecordEvent(obj, string(msg), "Set message to '%s'", string(msg))
	return nil
}

func (r *reconciler[S]) SetDeletingMessage(ctx context.Context, obj S, msg string) error {
	return r.SetConditionMessageByType(ctx, obj, ConditionReasonDeleting, msg)
}

func (r *reconciler[S]) SetReconcilingMessage(ctx context.Context, obj S, msg string) error {
	return r.SetConditionMessageByType(ctx, obj, ConditionTypeReady, msg)
}
