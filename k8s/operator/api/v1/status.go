package v1

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PhaseCode string

type ReconcilableStatus struct {
	State      PhaseCode          `json:"state"`
	Attempts   int32              `json:"attempts"`
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

func (r *ReconcilableStatus) SetState(state PhaseCode) {
	r.State = state
}

func (in *ReconcilableStatus) SetStatusCondition(condition metav1.Condition) {
	if in.Conditions == nil {
		in.Conditions = make([]metav1.Condition, 0)
	}
	conditions := in.Conditions
	meta.SetStatusCondition(
		&conditions,
		condition,
	)
	in.Conditions = conditions
}

func SetState(ctx context.Context, c client.Client, obj Resource, state string) error {
	log := log.FromContext(ctx)
	obj.ReconcilableStatus().State = PhaseCode(state)
	if err := c.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Failed to update Memcached status")
		return err
	}
	return nil
}

func (in *ReconcilableStatus) DeepCopy(out *ReconcilableStatus) {
	*out = *in
}

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}
