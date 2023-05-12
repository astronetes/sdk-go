package errors

var (
	CommunicationError = func(msg string) *ControllerError {
		return NewControllerError(communicationErrCode, msg)
	}
	CreateCloudResourceError = func(msg string) *ControllerError {
		return NewControllerError(createCloudResourceErrCode, msg)
	}
	DeleteCloudResourceError = func(msg string) *ControllerError {
		return NewControllerError(deleteCloudResourceErrCode, msg)
	}
	UnknownControllerError = func(msg string) *ControllerError {
		return NewControllerError(UnknownErrCode, msg)
	}
)

func NewControllerError(code ErrorCode, msg string) *ControllerError {
	return &ControllerError{
		ResourceError: NewResourceError(code, msg),
	}
}

type ControllerError struct {
	*ResourceError
	phase string
}

func (err *ControllerError) WithPhase(phase string) *ControllerError {
	err.phase = phase
	return err
}

func (err *ControllerError) Error() string {
	errMsg := ""
	if err.phase != "" {
		errMsg += "[ " + err.phase + " ] "
	}
	return errMsg + err.ResourceError.Error()
}
