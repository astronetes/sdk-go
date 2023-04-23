package v1

type Provider string

const (
	AWS   Provider = "aws"
	Azure Provider = "azure"
	K3s   Provider = "k3s"
	Kind  Provider = "kinds"
)
