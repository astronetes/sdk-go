package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PhaseCode string

const (
	FailedPhase      PhaseCode = "Failed"
	TerminatingPhase PhaseCode = "Terminating"
	ReadyPhase       PhaseCode = "Ready"
	DeletedPhase     PhaseCode = "Deleted"
)

/**
type Condition struct {
	metav1.Condition `json:",inline"`
	Causes           []string `json:"Causes,omitempty"`
}

type Conditions []Condition

func (c Conditions) isPreviousStatus(conditionType string) bool {
	return len(c) > 0 && conditionType == c[0].Type
}
**/

type ReconcilableStatus struct {
	Ready      bool               `json:"Ready"`
	State      PhaseCode          `json:"State"`
	Attempts   int32              `json:"Attempts"`
	Conditions []metav1.Condition `json:"Conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

func (in *ReconcilableStatus) SetReady(ready bool) {
	in.Ready = ready
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
	condition.Status = metav1.ConditionTrue
	if len(in.Conditions) == 0 {
		in.Conditions = []Condition{condition}
		return
	}
	if in.Conditions.isPreviousStatus(condition.Type) {
		in.Conditions[0].Message = condition.Message
		in.Conditions[0].LastTransitionTime = metav1.Now()
		return
	}

	in.Conditions[0].Status = metav1.ConditionFalse
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

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}
