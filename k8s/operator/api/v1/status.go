package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Condition struct {
	metav1.Condition `json:",inline"`
	Tries            int `json:"tries,omitempty"`
}

func NewConditionFromMetaV1(c metav1.Condition) Condition {
	return Condition{c, 0}
}

type Conditions []Condition

func (c Conditions) isPreviousStatus(conditionType string) bool {
	return !(len(c) > 0 && conditionType == c[0].Type)
}

type ReconcilableStatus struct {
	Ready         bool       `json:"Ready"`
	ErrorsCounter int        `json:"ErrorsCounter,omitempty"`
	Conditions    Conditions `json:"Conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=Conditions"`
	SpecHash      string     `json:"SpecHash,omitempty"`
	//+kubebuilder:validation:Optionals
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (s *ReconcilableStatus) ResetFailedCounter() {
	s.ErrorsCounter = 0
}

func (s *ReconcilableStatus) IncreaseEerrorCounter() {
	s.ErrorsCounter += 1
}

func (s *ReconcilableStatus) SetConditions(conditions ...Condition) {
	s.Conditions = conditions
}

func (s *ReconcilableStatus) updatePreviousState(condition Condition) {
	s.Conditions[0].Status = metav1.ConditionFalse
	s.Conditions[0].LastTransitionTime = metav1.Now()
	s.Conditions[0].Tries += 1
}

func (s *ReconcilableStatus) SetCondition(condition Condition) {
	conditions := s.Conditions
	conditions[0].Status = metav1.ConditionFalse
	conditions[0].LastTransitionTime = metav1.Now()
	// TODO move a to check
	exceedAllowedConditions := false
	startIndex := 0
	endIndex := len(conditions)
	if conditions.isPreviousStatus(condition.Type) {
		startIndex = 1
	} else if exceedAllowedConditions {
		endIndex -= 1
	}
	s.Conditions = append([]Condition{condition}, conditions[startIndex:endIndex]...)
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

func (r *ReconcilableStatus) ExceedErrors() bool {
	return r.ErrorsCounter > 3
}

func NewCondition(condType string, reason string, message string) Condition {
	return Condition{
		Condition: metav1.Condition{
			Type:               condType,
			Reason:             reason,
			Message:            message,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
		},
		Tries: 0,
	}
}
