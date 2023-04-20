package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var keyClient = "client"

func AddCLientToContext(parent context.Context, client client.Client) context.Context {
	return context.WithValue(parent, keyClient, client)
}

func ClientFromContext(ctx context.Context) client.Client {
	value := ctx.Value(keyClient)
	if value == nil {
		return nil
	}
	return value.(client.Client)
}
