package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PhaseCode string

type (
	OperationID     string
	OperationStatus string
)

const (
	Creation   OperationID     = "create"
	Updation   OperationID     = "update"
	Deletion   OperationID     = "delete"
	OnCreation OperationStatus = "onCreation"
	Ready      OperationStatus = "ready"
	OnDeletion OperationStatus = "onDeletion"
	Deleted    OperationStatus = "deleted"
)

type Operation struct {
	Op         OperationID `json:"op,omitempty"`
	ResourceID string      `json:"resourceID,omitempty"`
}

type Condition struct {
	metav1.Condition `json:",inline"`
	Tries            int         `json:"tries,omitempty"`
	Operations       []Operation `json:"operations,omitempty"`
}

type Conditions []Condition

func (c Conditions) isPreviousStatus(conditionType string) bool {
	return len(c) > 0 && conditionType == c[0].Type
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

func (s *ReconcilableStatus) AddCondition(condition Condition) {
	conditions := s.Conditions
	conditions[0].Status = metav1.ConditionFalse
	conditions[0].LastTransitionTime = metav1.Now()
	// TODO move a to check
	exceedAllowedConditions := false
	startIndex := 0
	endIndex := len(conditions)
	if conditions.isPreviousStatus(condition.Type) {
		startIndex = 1
		condition.Tries = conditions[0].Tries + 1
	} else if exceedAllowedConditions {
		endIndex -= 1
	}
	s.Conditions = append([]Condition{condition}, conditions[startIndex:endIndex]...)
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

func (r *ReconcilableStatus) ExceedErrors() bool {
	return r.ErrorsCounter > 3
}

func NewCondition(condType PhaseCode, reason string, message string) Condition {
	return Condition{
		Condition: metav1.Condition{
			Type:               string(condType),
			Reason:             reason,
			Message:            message,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
		},
		Tries: 0,
	}
}
