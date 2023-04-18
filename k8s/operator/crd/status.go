package crd

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Status struct {
	Ready      bool               `json:"ready"`
	Failed     int                `json:"failed,omitempty"`
	Phase      metav1.Condition   `json:"phase,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	SpecHash   string             `json:"specHash,omitempty"`
	//+kubebuilder:validation:Optional
	ErrorMessage string `json:"errorMessage,omitempty"`
}
