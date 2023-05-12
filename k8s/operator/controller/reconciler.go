package controller

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Astronetes[S v1.Resource] interface {
	Reconcile(ctx context.Context, req ctrl.Request, obj S) (ctrl.Result, error)
}

type astronetes[S v1.Resource] struct {
	client.Client
	ID            string
	FinalizerName string
	Recorder      record.EventRecorder
	Tracer        trace.Tracer
	Config        config.Controller
	Scheme        *runtime.Scheme
	Dispatcher    Dispatcher[S]
	corePhases    map[v1.PhaseCode]PhaseReconcile[S]
}

func (r *astronetes[S]) loadCorePhases() {
	if r.corePhases == nil {
		r.corePhases = make(map[v1.PhaseCode]PhaseReconcile[S])
	}
	r.corePhases[v1.ReadyPhase] = r.reconcileReady
	r.corePhases[v1.FailedPhase] = r.reconcileFailed
	r.corePhases[v1.TerminatingPhase] = r.reconcileTerminating
	r.corePhases[v1.DeletedPhase] = r.reconcileDeleted
}

func NewAstronetesReconcile[S v1.Resource](id string, mgr manager.Manager, finalizerName string,
	config config.Controller, dispatcher Dispatcher[S],
) Astronetes[S] {
	a := &astronetes[S]{
		Client:        mgr.GetClient(),
		ID:            id,
		FinalizerName: finalizerName,
		Recorder:      mgr.GetEventRecorderFor(id),
		Tracer:        otel.Tracer(id),
		Config:        config,
		Scheme:        mgr.GetScheme(),
		Dispatcher:    dispatcher,
	}
	a.loadCorePhases()
	return a
}

func (r *astronetes[S]) Reconcile(ctx context.Context, req ctrl.Request, obj S) (ctrl.Result, error) {
	// Initialize the logger
	log := log.FromContext(ctx)
	log.Info("reconciling Ingress Controller")
	// Initialize the Tracer
	ctx, span := r.Tracer.Start(ctx, fmt.Sprintf("%s-reconcile", r.ID))
	defer span.End()
	status := obj.AstronetesStatus()
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	currentPhaseCode := v1.PhaseCode(status.GetCurrentPhase())

	span.SetAttributes(attribute.String("phase", string(currentPhaseCode)))
	span.SetAttributes(attribute.Int("attempt", int(status.Attempts)))
	span.SetAttributes(attribute.String("resource", obj.GetName()))
	cfg := r.Config.GetConfigForReconciliationPhase(currentPhaseCode)

	// Add a deadline just to make sure we don't get stuck in a loop
	ctx, cancel := context.WithDeadline(ctx, metav1.Now().Add(*cfg.Timeout))
	defer cancel()

	reachLimitOfTries := cfg.AllowedAttempts > 0 && status.Attempts > cfg.AllowedAttempts

	// The maximum number of attempts for the current phase is reached
	if reachLimitOfTries {
		span.AddEvent(fmt.Sprintf("Attempt number %v", status.Attempts))
		err := fmt.Errorf("reach max number of allowed attempts (%v) for completing the phase '%s'", status.Attempts, currentPhaseCode)
		status.AddErrorCause(err)
		status.Next(v1.FailedPhase, ReachMaxAllowedAttemptsEvent, "phase not completed in max number of allowed attempts")
		if err := r.Status().Update(ctx, obj); err != nil {
			log.Info(err.Error())
			return ctrl.Result{}, err
		}
		r.Recorder.Event(obj, corev1.EventTypeWarning, ReachMaxAllowedAttemptsEvent, err.Error())
		return ctrl.Result{Requeue: true}, nil
	}

	// Shouldn't be checked other possible flows If the status hasn't been initialized yet.
	if "" == currentPhaseCode {
		status.Next(r.Dispatcher.InitialCreationPhaseCode, NewRequestEvent, "Starting the creation of resources")
		status.Ready = false
		status.State = v1.PhaseCode(status.Conditions[0].Type)
		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Info(err.Error())
		}
		r.Recorder.Eventf(obj, corev1.EventTypeNormal, NewRequestEvent, "Processing the request for creating the resource")
		return ctrl.Result{Requeue: true}, nil
	}

	// This is a special check to verify if the resource is already in "in deletion" (temporary status)
	if !obj.GetDeletionTimestamp().IsZero() && !r.IsOnDeletionPhase(currentPhaseCode) {
		status.Next(v1.TerminatingPhase, TerminatingEvent, "Terminating the resource")
		status.Ready = false
		status.State = v1.PhaseCode(status.Conditions[0].Type)
		r.Recorder.Event(obj, corev1.EventTypeNormal, "Terminating resource", "Setting status as 'Terminating' for the resource")
		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Info(err.Error())
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// This is a special check to verify if we need to add the finalizer
	if !controllerutil.ContainsFinalizer(obj, r.FinalizerName) {
		controllerutil.AddFinalizer(obj, r.FinalizerName)
		span.AddEvent(fmt.Sprintf("Adding finalizer '%s' to resource", r.FinalizerName))
		r.Recorder.Eventf(obj, corev1.EventTypeNormal, "Add finalizer", "Adding finalizer '%s' to resource", r.FinalizerName)
		if err := r.Client.Update(ctx, obj); err != nil {
			log.Info(err.Error())
		}
		return ctrl.Result{Requeue: true}, nil
	}

	status.Attempts += 1
	result, err := r.ReconcilePhase(ctx, currentPhaseCode, r.Client, cfg, obj)
	if err != nil {
		span.RecordError(err)
		return ctrl.Result{Requeue: false}, err
	}
	status.State = v1.PhaseCode(status.Conditions[0].Type)
	if err := r.Client.Status().Update(ctx, obj); err != nil {
		log.Info(err.Error())
	}
	switch result.Code() {
	case ErrorCode:
		r.Recorder.Event(obj, corev1.EventTypeWarning, ErrorEvent, result.Message())
		if requeueAfter := cfg.GetRequeueAfterByAttemptNumber(status.Attempts); requeueAfter != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: *requeueAfter,
			}, nil
		}
		return ctrl.Result{}, fmt.Errorf(result.Message())
	case OKCode:
		r.Recorder.Event(obj, corev1.EventTypeNormal, status.Conditions[0].Reason, result.Message())
		return ctrl.Result{
			Requeue: true,
		}, nil
	case CompletedCode:
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, fmt.Errorf("unexpected status")
}

func (r *astronetes[S]) ReconcilePhase(ctx context.Context, code v1.PhaseCode, c client.Client, cfg config.Phase, obj S) (Result, error) {
	p, ok := r.corePhases[code]
	if !ok {
		return r.Dispatcher.ReconcilePhase(ctx, code, c, cfg, obj)
	}
	return p(ctx, c, cfg, obj), nil
}

func (r *astronetes[S]) IsOnDeletionPhase(code v1.PhaseCode) bool {
	return code == v1.TerminatingPhase || code == v1.DeletedPhase || r.Dispatcher.IsOnDeletionPhase(code)
}
