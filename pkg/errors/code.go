package errors

type ErrorCode string

const (
	UnknownErrCode ErrorCode = "ErrUnknown"
)

const (
	InvalidRequestErrCode           ErrorCode = "ErrInvalidRequest"
	MissingRequiredAttributeErrCode ErrorCode = "ErrMissingRequiredAttribute"
)

// TODO Should we use different error codes for Resource and Operator errors?
const (
	CommunicationErrCode       ErrorCode = "ErrCommunication"
	CreateCloudResourceErrCode ErrorCode = "ErrCreateCloudResource"
	DeleteCloudResourceErrCode ErrorCode = "ErrDeleteCloudResource"
)
