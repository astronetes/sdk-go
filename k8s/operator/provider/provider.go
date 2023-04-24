package provider

import "context"

type Provider[T any] interface {
	Execute(ctx context.Context, obj T) error
	IsReady(ctx context.Context, obj T) (bool, error)
}
