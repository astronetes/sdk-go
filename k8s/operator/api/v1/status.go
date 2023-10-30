package v1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PhaseCode string

type ReconcilableStatus struct {
	State      PhaseCode          `json:"state"`
	Attempts   int32              `json:"attempts"`
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" 
		patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

func (r *ReconcilableStatus) SetState(state PhaseCode) {
	r.State = state
}

func (in *ReconcilableStatus) GetStatusCondition(conditionType string) *metav1.Condition {
	if in.Conditions == nil {
		return nil
	}

	for _, condition := range in.Conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}

	return nil
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

func (in *ReconcilableStatus) DeepCopy(out *ReconcilableStatus) {
	*out = *in
}

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}
