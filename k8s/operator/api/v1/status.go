package v1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PhaseCode string

type ReconcilableStatus struct {
	State      PhaseCode          `json:"state"`
	Attempts   int32              `json:"attempts"`
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
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

/**
func (in *ReconcilableStatus) Next(phase PhaseCode, event string, msg string) {
	in.Attempts = 0
	condition := Condition{
		Condition: metav1.Condition{
			Type:               string(phase),
			Reason:             event,
			Message:            msg,
			LastTransitionTime: metav1.Now(),
		},
	}
	in.addCondition(condition)
}
**/
/*
*

	func (s *ReconcilableStatus) updatePreviousState(condition Condition) {
		s.Conditions[0].LastTransitionTime = metav1.Now()
		s.Conditions[0].Message = condition.Message
		s.Conditions[0].Attempts += 1
	}

*
*/
/**
func (in *ReconcilableStatus) addCondition(condition Condition) {
	condition.AstronetesStatus = metav1.ConditionTrue
	if len(in.Conditions) == 0 {
		in.Conditions = []Condition{condition}
		return
	}
	if in.Conditions.isPreviousStatus(condition.Type) {
		in.Conditions[0].Message = condition.Message
		in.Conditions[0].LastTransitionTime = metav1.Now()
		return
	}

	in.Conditions[0].AstronetesStatus = metav1.ConditionFalse
	exceedAllowedConditions := len(in.Conditions) > 10

	endIndex := len(in.Conditions)
	if exceedAllowedConditions {
		endIndex -= 1
	}
	in.Conditions = append([]Condition{condition}, in.Conditions[0:endIndex]...)
}

func (in *ReconcilableStatus) GetCurrentCondition() Condition {
	if len(in.Conditions) > 0 {
		return in.Conditions[0]
	}
	return Condition{}
}

func (in *ReconcilableStatus) GetCurrentPhase() string {
	if len(in.Conditions) > 0 {
		return in.Conditions[0].Type
	}
	return ""
}

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}

/**
func (in *ReconcilableStatus) AddErrorCause(err error) {
	if in.Conditions[0].Causes == nil {
		in.Conditions[0].Causes = make([]string, 0)
	}
	in.Conditions[0].Causes = append(in.Conditions[0].Causes, err.Error())
}
*/
