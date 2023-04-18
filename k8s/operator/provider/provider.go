package provider

type Provider int

const (
	AWS Provider = iota
	Kind
	K3s
	GCP
	Azure
)
