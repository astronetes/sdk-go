package errors

import "fmt"

var (
	MissingRequiredAttributeError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(missingRequiredAttributeErrCode, msg, args...)
	}
	InvalidRequestError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(invalidRequestErrCode, msg, args...)
	}
	UnknownResourceError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(invalidRequestErrCode, msg, args...)
	}
)

func NewResourceError(code ErrorCode, msg string, args ...any) *ResourceError {
	return &ResourceError{
		AstronetesError: AstronetesError{
			code: code,
			msg:  fmt.Sprintf(msg, args),
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
