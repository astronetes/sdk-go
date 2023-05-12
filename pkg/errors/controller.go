package errors

var (
	CommunicationError = func(msg string) *ControllerError {
		return &ControllerError{
			ResourceError: NewResourceError(communicationErrCode, msg),
		}
	}
	CreateCloudResourceError = func(msg string) *ControllerError {
		return &ControllerError{
			ResourceError: NewResourceError(createCloudResourceErrCode, msg),
		}
	}
	DeleteCloudResourceError = func(msg string) *ControllerError {
		return &ControllerError{
			ResourceError: NewResourceError(deleteCloudResourceErrCode, msg),
		}
	}
	UnknownControllerError = func(msg string) *ResourceError {
		return &ResourceError{
			AstronetesError: AstronetesError{
				code: UnknownErrCode,
				msg:  msg,
			},
		}
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
