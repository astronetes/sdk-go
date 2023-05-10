package reconcile

import (
	"context"
	"fmt"

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

type Astronetes struct {
	client.Client
	ID            string
	FinalizerName string
	Recorder      record.EventRecorder
	Tracer        trace.Tracer
	Config        config.Controller
	Scheme        *runtime.Scheme
	Dispatcher    Dispatcher
}

func (r *Astronetes) Reconcile(ctx context.Context, req ctrl.Request, obj v1.Resource) (ctrl.Result, error) {
	// Initialize the logger
	log := log.FromContext(ctx)
	log.Info("reconciling Ingress Controller")
	// Initialize the Tracer
	ctx, span := r.Tracer.Start(ctx, fmt.Sprintf("%s-reconcile", r.ID))
	defer span.End()

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	currentPhaseCode := v1.PhaseCode(obj.AstronetesStatus().GetCurrentPhase())

	span.SetAttributes(attribute.String("phase", string(currentPhaseCode)))
	span.SetAttributes(attribute.Int("attempt", int(obj.AstronetesStatus().Attempts)))
	span.SetAttributes(attribute.String("resource", obj.GetName()))
	cfg := r.Config.GetConfigForReconciliationPhase(currentPhaseCode)

	// Add a deadline just to make sure we don't get stuck in a loop
	ctx, cancel := context.WithDeadline(ctx, metav1.Now().Add(*cfg.Timeout))
	defer cancel()

	reachLimitOfTries := cfg.AllowedAttempts > 0 && obj.AstronetesStatus().Attempts > cfg.AllowedAttempts

	// The maximum number of attempts for the current phase is reached
	if reachLimitOfTries {
		span.AddEvent(fmt.Sprintf("Attempt number %v", obj.AstronetesStatus().Attempts))
		err := fmt.Errorf("reach max number of allowed attempts (%v) for completing the phase '%s'", obj.AstronetesStatus().Attempts, currentPhaseCode)
		obj.AstronetesStatus().AddErrorCause(err)
		obj.AstronetesStatus().Next(v1.FailedPhase, ReachMaxAllowedAttemptsEvent, "phase not completed in max number of allowed attempts")
		if err := r.Status().Update(ctx, obj); err != nil {
			log.Info(err.Error())
			return ctrl.Result{}, err
		}
		r.Recorder.Event(obj, corev1.EventTypeWarning, ReachMaxAllowedAttemptsEvent, err.Error())
		return ctrl.Result{Requeue: true}, nil
	}

	// Shouldn't be checked other possible flows If the status hasn't been initialized yet.
	if "" == currentPhaseCode {
		obj.AstronetesStatus().Next(r.Dispatcher.InitialPhaseCode, NewRequestEvent, "Starting the creation of resources")
		obj.AstronetesStatus().Ready = false
		obj.AstronetesStatus().State = v1.PhaseCode(obj.AstronetesStatus().Conditions[0].Type)
		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Info(err.Error())
		}
		r.Recorder.Eventf(obj, corev1.EventTypeNormal, NewRequestEvent, "Processing the request for creating the resource")
		return ctrl.Result{Requeue: true}, nil
	}

	// This is a special check to verify if the resource is already in "in deletion" (temporary status)
	if !obj.GetDeletionTimestamp().IsZero() && !r.Dispatcher.IsOnDeletionPhase(currentPhaseCode) {
		obj.AstronetesStatus().Next(v1.TerminatingPhase, TerminatingEvent, "Terminating the resource")
		obj.AstronetesStatus().Ready = false
		obj.AstronetesStatus().State = v1.PhaseCode(obj.AstronetesStatus().Conditions[0].Type)
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

	obj.AstronetesStatus().Attempts += 1
	reconcileFn, err := r.Dispatcher.GetPhase(currentPhaseCode)
	if err != nil {
		span.RecordError(err)
		return ctrl.Result{Requeue: false}, err
	}
	result := reconcileFn(ctx, r.Client, cfg, obj)
	obj.AstronetesStatus().State = v1.PhaseCode(obj.AstronetesStatus().Conditions[0].Type)
	if err := r.Client.Status().Update(ctx, obj); err != nil {
		log.Info(err.Error())
	}
	switch result.Code() {
	case ErrorCode:
		r.Recorder.Event(obj, corev1.EventTypeWarning, ErrorEvent, result.Message())
		if requeueAfter := cfg.GetRequeueAfterByAttemptNumber(obj.AstronetesStatus().Attempts); requeueAfter != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: *requeueAfter,
			}, nil
		}
		return ctrl.Result{}, fmt.Errorf(result.Message())
	case OKCode:
		r.Recorder.Event(obj, corev1.EventTypeNormal, obj.AstronetesStatus().Conditions[0].Reason, result.Message())
		return ctrl.Result{
			Requeue: true,
		}, nil
	case CompletedCode:
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, fmt.Errorf("unexpected status")
}
