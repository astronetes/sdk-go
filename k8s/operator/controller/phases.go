package controller

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

type PhaseCode string

type PhaseResultCode int32

type PhaseResult struct {
	Code           PhaseResultCode
	Msg            string
	Error          error
	ShouldRetry    bool
	ReconcileAfter time.Duration
}

func (r *PhaseResult) HasError() bool {
	return r.Error != nil
}

func (r *PhaseResult) RuntimeResult() (ctrl.Result, error) {
	switch r.Code {
	case Ok, Retriable:
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: r.ReconcileAfter,
		}, r.Error
	case Error:
		return ctrl.Result{}, r.Error
	case Completed:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, nil
	}
}

const (
	Ok        PhaseResultCode = iota
	Retriable PhaseResultCode = iota
	Completed PhaseResultCode = iota
	Error     PhaseResultCode = iota
)

var ResultOK = func(msg string) PhaseResult {
	return PhaseResult{
		Code: Ok,
		Msg:  msg,
	}
}

var ResultRetriable = func(msg string, err error, nextRetryIn time.Duration) PhaseResult {
	return PhaseResult{
		Code:           Retriable,
		Msg:            msg,
		Error:          err,
		ReconcileAfter: nextRetryIn,
	}
}

var ResultError = func(err error) PhaseResult {
	return PhaseResult{
		Code:  Error,
		Error: err,
	}
}

var ResultReconciliationCompleted = func(err error) PhaseResult {
	return PhaseResult{
		Code: Completed,
	}
}

type PhaseReconciler interface {
	Reconcile(ctx context.Context) PhaseResult
}
