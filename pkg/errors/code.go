package errors

type ErrorCode string

const (
	UnknownErrCode ErrorCode = "ErrUnknown"
)

const (
	invalidRequestErrCode           ErrorCode = "ErrInvalidRequest"
	missingRequiredAttributeErrCode ErrorCode = "ErrMissingRequiredAttribute"
)

// TODO Should we use different error codes for Resource and Operator errors?
const (
	communicationErrCode       ErrorCode = "ErrCommunication"
	createCloudResourceErrCode ErrorCode = "ErrCreateCloudResource"
	deleteCloudResourceErrCode ErrorCode = "ErrDeleteCloudResource"
)
