package errors

import "fmt"

var (
	MissingRequiredAttributeError = func(msg string) *ResourceError {
		return &ResourceError{
			AstronetesError: AstronetesError{
				code: missingRequiredAttributeErrCode,
				msg:  msg,
			},
		}
	}
	InvalidRequestError = func(msg string) *ResourceError {
		return &ResourceError{
			AstronetesError: AstronetesError{
				code: invalidRequestErrCode,
				msg:  msg,
			},
		}
	}
	UnknownResourceError = func(msg string) *ResourceError {
		return &ResourceError{
			AstronetesError: AstronetesError{
				code: invalidRequestErrCode,
				msg:  msg,
			},
		}
	}
)

func NewResourceError(code ErrorCode, msg string) *ResourceError {
	return &ResourceError{
		AstronetesError: AstronetesError{
			code: code,
			msg:  msg,
		},
	}
}

type ResourceError struct {
	AstronetesError
	ns       string
	resource string
	name     string
}

func (err *ResourceError) WithNS(ns string) *ResourceError {
	err.ns = ns
	return err
}

func (err *ResourceError) WithResource(resource string) *ResourceError {
	err.resource = resource
	return err
}

func (err *ResourceError) WithName(name string) *ResourceError {
	err.name = name
	return err
}

func (err *ResourceError) Error() string {
	resourcePath := ""
	if err.ns != "" || err.resource != "" {
		resourcePath = err.ns + " / " + err.resource
	}
	if err.name != "" {
		resourcePath += "( " + err.name + " )"
	}
	return fmt.Sprintf("%s [%v] %s", resourcePath, err.code, err.msg)
}
