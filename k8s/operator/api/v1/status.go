package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PhaseCode string

type Condition struct {
	metav1.Condition `json:",inline"`
	// Attempts         int32 `json:"attempts,omitempty"`
}

type Conditions []Condition

func (c Conditions) isPreviousStatus(conditionType string) bool {
	return len(c) > 0 && conditionType == c[0].Type
}

type ReconcilableStatus struct {
	Ready      bool       `json:"Ready"`
	State      PhaseCode  `json:"state"`
	Attempts   int32      `json:"Attempts"`
	Conditions Conditions `json:"Conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

func (r *ReconcilableStatus) Next(phase PhaseCode, event string, msg string) {
	r.Attempts = 0
	condition := Condition{
		Condition: metav1.Condition{
			Type:               string(phase),
			Reason:             event,
			Message:            msg,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
		},
	}
	r.addCondition(condition)
}

/*
*

	func (s *ReconcilableStatus) updatePreviousState(condition Condition) {
		s.Conditions[0].LastTransitionTime = metav1.Now()
		s.Conditions[0].Message = condition.Message
		s.Conditions[0].Attempts += 1
	}

*
*/
func (s *ReconcilableStatus) addCondition(condition Condition) {
	if len(s.Conditions) == 0 {
		s.Conditions = []Condition{condition}
		return
	}
	conditions := s.Conditions
	conditions[0].Status = metav1.ConditionFalse
	conditions[0].LastTransitionTime = metav1.Now()
	// TODO move a to check
	exceedAllowedConditions := len(conditions) > 10
	if conditions.isPreviousStatus(condition.Type) {
		s.Conditions[0].Message = condition.Message
		// s.updatePreviousState(condition)
		return
	}
	endIndex := len(conditions)
	if exceedAllowedConditions {
		endIndex -= 1
	}
	s.Conditions = append([]Condition{condition}, conditions[0:endIndex]...)
	s.Conditions[0].Status = metav1.ConditionTrue
}

func (s *ReconcilableStatus) GetCurrentCondition() Condition {
	if len(s.Conditions) > 0 {
		return s.Conditions[0]
	}
	return Condition{}
}

func (s *ReconcilableStatus) GetCurrentPhase() string {
	if len(s.Conditions) > 0 {
		return s.Conditions[0].Type
	}
	return ""
}

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}

/**
func NewCondition(condType PhaseCode, reason string, message string) Condition {
	return Condition{
		Condition: metav1.Condition{
			Type:               string(condType),
			Reason:             reason,
			Message:            message,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
		},
		// Attempts: 1,
	}
}

*/
