package provider

import "context"

type OperationResultCode int32

const (
	Ok OperationResultCode = iota
	Error
)

type Provider[T any] interface {
	Create(ctx context.Context, res T) (OperationResult, error)
	Destroy(ctx context.Context, res T) (OperationResult, error)
	Reconcile(ctx context.Context, res T) (OperationResult, error)
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
