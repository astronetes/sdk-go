package errors

type ErrorCode string

const (
	InvalidRequestErrCode           ErrorCode = "ErrInvalidRequest"
	MissingRequiredAttributeErrCode ErrorCode = "ErrMissingRequiredAttribute"
)
