package controller

import (
	"context"
	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"github.com/astronetes/sdk-go/k8s/operator/config"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	OKCode Code = iota
	CompletedCode
	ErrorCode
)

var OK = func(msg string) Result {
	return Result{
		code: OKCode,
		msg:  msg,
	}
}

var Error = func(err error) Result {
	return Result{
		code: ErrorCode,
		err:  err,
	}
}

var Completed = func(msg string) Result {
	return Result{
		code: CompletedCode,
		msg:  msg,
	}
}

type PhaseReconcile[O v1.Resource] func(ctx context.Context, client client.Client, cfg config.Phase, obj O) Result

type Result struct {
	code  Code
	msg   string
	err   error
	after *time.Duration
}

type Code int32

func (r Result) HasError() bool {
	return r.err != nil
}

func (r Result) Message() string {
	if r.err != nil {
		return r.err.Error()
	}
	return r.msg
}

func (r Result) Code() Code {
	return r.code
}

func (r Result) After(t *time.Duration) Result {
	r.after = t
	return r
}

func (r Result) RuntimeResult() (ctrl.Result, error) {
	shouldRequeue := r.code != CompletedCode && (r.Code() == OKCode || r.after != nil)
	if !shouldRequeue {
		return ctrl.Result{}, r.err
	}
	var requeueAfter time.Duration
	if r.after != nil {
		requeueAfter = *r.after
	}
	return ctrl.Result{
		RequeueAfter: requeueAfter,
	}, nil
}
