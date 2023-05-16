package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *reconciler[S]) addFinalizer(ctx context.Context, req ctrl.Request) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	var obj S

	// Fetch the latest Memcached
	// If this fails, bubble up the reconcile results to the main reconciler
	if r, err := r.getLatest(ctx, req, obj); ShouldHaltOrRequeue(r, err) {
		return r, err
	}
	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(obj, r.finalizerName) {
		log.Info("Adding Finalizer for the resource")
		if ok := controllerutil.AddFinalizer(obj, r.finalizerName); !ok {
			log.Error(nil, "Failed to add finalizer into the custom resource")
			return Requeue()
		}

		if err := r.Update(ctx, obj); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return RequeueWithError(err)
		}
	}

	return ContinueReconciling()
}
