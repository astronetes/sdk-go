package provider

import "context"

type OperationResultCode int32

const (
	Ok OperationResultCode = iota
	Error
)

type Provider[T any] interface {
	Create(ctx context.Context, res T) (OperationResult, error)
	IsReady(ctx context.Context, res T) (bool, error)
	Destroy(ctx context.Context, res T) (OperationResult, error)
	IsDestroyed(ctx context.Context, res T) (bool, error)
}

type OperationResult interface {
	Code() OperationResultCode
}

type GenericOperationResult struct {
	code OperationResultCode
}

func (r *GenericOperationResult) Code() OperationResultCode {
	return r.code
}
