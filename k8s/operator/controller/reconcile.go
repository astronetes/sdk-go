package controller

import (
	"context"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	okCode Code = iota
	doneCode
	failCode
)

var OK = func(msg string) Result {
	return Result{
		code: okCode,
		msg:  msg,
	}
}

var Error = func(err error) Result {
	return Result{
		code: failCode,
		err:  err,
	}
}

var Completed = func(msg string) Result {
	return Result{
		code: doneCode,
		msg:  msg,
	}
}

type PhaseReconcile[O any] func(ctx context.Context, cfg config.Phase, obj O) Result

type Result struct {
	code  Code
	msg   string
	err   error
	after time.Duration
}

type Code int32

func (r Result) HasError() bool {
	return r.err != nil
}

func (r Result) Message() string {
	return r.msg
}

func (r Result) Code() Code {
	return r.code
}

func (r Result) After(t time.Duration) Result {
	r.after = t
	return r
}

func (r Result) RuntimeResult() (ctrl.Result, error) {
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
