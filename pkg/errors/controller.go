package errors

var (
	CommunicationError = func(msg string, args ...any) *ControllerError {
		return NewControllerError(CommunicationErrCode, msg, args)
	}
	CreateCloudResourceError = func(msg string, args ...any) *ControllerError {
		return NewControllerError(CreateCloudResourceErrCode, msg, args)
	}
	DeleteCloudResourceError = func(msg string, args ...any) *ControllerError {
		return NewControllerError(DeleteCloudResourceErrCode, msg, args)
	}
	CloudResourceNotFoundError = func(msg string, args ...any) *ControllerError {
		return NewControllerError(CloudResourceNotFoundErrCode, msg, args)
	}
	UnknownControllerError = func(msg string, args ...any) *ControllerError {
		return NewControllerError(UnknownErrCode, msg, args)
	}
)

func NewControllerError(code ErrorCode, msg string, args ...any) *ControllerError {
	return &ControllerError{
		ResourceError: NewResourceError(code, msg, args...),
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

func (err *ControllerError) Wrap(er error) *ControllerError {
	err.err = er
	return err
}

func (err *ControllerError) Error() string {
	errMsg := ""
	if err.phase != "" {
		errMsg += "[ " + err.phase + " ] "
	}
	return errMsg + err.ResourceError.Error()
}
