package reconciler

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// FnWithRequest is a function definition representing small
// reconciliation behavior. The request is included as a parameter.
type FnWithRequest = func(ctx context.Context, c client.Client, cfg Config, req ctrl.Request) (*ctrl.Result, error)

// Requeue returns a controller result pairing specifying to
// requeue with no error message implied. This returns no error.
func Requeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, nil }

// DoNotRequeue returns a controller result pairing specifying not to requeue.
func DoNotRequeue() (*reconcile.Result, error) { return &ctrl.Result{Requeue: false}, nil }

// RequeueWithError returns a controller result pairing specifying to
// requeue with an error message.
func RequeueWithError(e error) (*reconcile.Result, error) { return &ctrl.Result{Requeue: true}, e }

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

func (r *reconciler[S]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// The list of subreconcilers for Memcached
	subreconcilersForResource := []FnWithRequest{
		r.startReconciliation,
		r.addFinalizer,
		r.handleDeletion,
		r.handleCreation,
		r.updateStatus,
	}

	// Run all subreconcilers sequentially
	for _, f := range subreconcilersForResource {
		if r, err := f(ctx, r.Client, r.config, req); ShouldHaltOrRequeue(r, err) {
			return Evaluate(r, err)
		}
	}

	return Evaluate(DoNotRequeue())
}
