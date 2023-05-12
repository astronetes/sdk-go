package errors

type ErrorCode string

const (
	invalidRequestErrCode           ErrorCode = "ErrInvalidRequest"
	missingRequiredAttributeErrCode ErrorCode = "ErrMissingRequiredAttribute"
)
