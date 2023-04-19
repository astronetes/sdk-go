package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Condition struct {
	metav1.Condition `json:",inline"`
}

func NewConditionFromMetaV1(c metav1.Condition) Condition {
	return Condition{c}
}

type Conditions []Condition

func (c Conditions) isPreviousStatus(conditionType string) bool {
	return !(len(c) > 0 && conditionType == c[0].Type)
}

type ReconcilableStatus struct {
	ready         bool       `json:"ready"`
	errorsCounter int        `json:"errorsCounter,omitempty"`
	conditions    Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	specHash      string     `json:"specHash,omitempty"`
	//+kubebuilder:validation:Optionals
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (s ReconcilableStatus) SetReady(ready bool) {
	s.ready = ready
}

func (s ReconcilableStatus) ResetFailedCounter() {
	s.errorsCounter = 0
}

func (s ReconcilableStatus) IncreaseEerrorCounter() {
	s.errorsCounter += 1
}

func (s ReconcilableStatus) SetConditions(conditions ...Condition) {
	s.conditions = conditions
}

func (s ReconcilableStatus) SetCondition(condition Condition) {
	conditions := s.conditions
	conditions[0].Status = metav1.ConditionFalse
	exceedAllowedConditions := false
	startIndex := 0
	endIndex := len(conditions)
	if conditions.isPreviousStatus(condition.Type) {
		startIndex = 1
	} else if exceedAllowedConditions {
		endIndex -= 1
	}
	s.conditions = append([]Condition{condition}, conditions[startIndex:endIndex]...)
}

func (s ReconcilableStatus) Conditions() Conditions {
	return s.conditions
}

func (s ReconcilableStatus) SetSpecHash(state string) {
	s.specHash = state
}

func (s ReconcilableStatus) GetCurrentPhase() string {
	if len(s.conditions) > 0 {
		return s.conditions[0].Type
	}
	return ""
}

func (in *ReconcilableStatus) DeepCopyInto(out *ReconcilableStatus) {
	*out = *in
}

func (r *ReconcilableStatus) SpecHash() string {
	return r.specHash
}

func (r *ReconcilableStatus) ExceedErrors() bool {
	return r.errorsCounter > 3
}
