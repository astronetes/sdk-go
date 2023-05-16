package errors

var (
	CommunicationError = func(msg string, args ...any) error {
		return newControllerError(CommunicationErrCode, msg, args)
	}
	CreateCloudResourceError = func(msg string, args ...any) error {
		return newControllerError(CreateCloudResourceErrCode, msg, args)
	}
	DeleteCloudResourceError = func(msg string, args ...any) error {
		return newControllerError(DeleteCloudResourceErrCode, msg, args)
	}
	CloudResourceNotFoundError = func(msg string, args ...any) error {
		return newControllerError(CloudResourceNotFoundErrCode, msg, args)
	}
	UnknownControllerError = func(msg string, args ...any) error {
		return newControllerError(UnknownErrCode, msg, args)
	}
)

func newControllerError(code ErrorCode, msg string, args ...any) *ControllerError {
	return &ControllerError{
		ResourceError: newResourceError(code, msg, args...),
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
