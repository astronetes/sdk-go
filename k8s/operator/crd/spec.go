package crd

type Spec struct {

	// ClassName to be assigned to the Controller
	ClassName string `json:"classname,omitempty"`

	// Set the Controller as default one
	//+kubebuilder:default:=false
	Default bool `json:"default,omitempty"`

	//+kubebuilder:validation:Optional
	provider Provider `json:"provider,omitempty"`
}
