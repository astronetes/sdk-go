package controller

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	okCode PhaseResultCode = iota
	doneCode
	failCode
)

var ReconciliationOk = func(msg string) PhaseResult {
	return PhaseResult{
		code: okCode,
		msg:  msg,
	}
}

var ReconciliationError = func(err error) PhaseResult {
	return PhaseResult{
		code: failCode,
		err:  err,
	}
}

var ReconciliationCompleted = func(msg string) PhaseResult {
	return PhaseResult{
		code: doneCode,
		msg:  msg,
	}
}

type PhaseReconcile[O any] func(ctx context.Context, obj O) PhaseResult

type PhaseResult struct {
	code  PhaseResultCode
	msg   string
	err   error
	after time.Duration
}

type PhaseResultCode int32

func (r PhaseResult) HasError() bool {
	return r.err != nil
}

func (r PhaseResult) Message() string {
	return r.msg
}

func (r PhaseResult) Code() PhaseResultCode {
	return r.code
}

func (r PhaseResult) After(t time.Duration) PhaseResult {
	r.after = t
	return r
}

func (r *PhaseResult) RuntimeResult() (ctrl.Result, error) {
	switch r.code {
	case okCode:
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: r.after,
		}, r.err
	case failCode:
		return ctrl.Result{}, r.err
	case doneCode:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, r.err
	}
}
