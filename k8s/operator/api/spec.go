package api

type Provider string

const (
	AWS   Provider = "aws"
	Azure Provider = "azure"
	K3s   Provider = "k3s"
	Kind  Provider = "kinds"
)

type Spec struct {
	// ClassName to be assigned to the Controller
	ClassName string `json:"classname,omitempty"`

	// Set the Controller as default one
	//+kubebuilder:default:=false
	Default bool `json:"default,omitempty"`

	//+kubebuilder:validation:Optional
	Provider Provider `json:"Provider,omitempty"`
}
