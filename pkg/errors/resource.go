package errors

import "fmt"

var (
	MissingRequiredAttributeError = func(msg string) *ResourceError {
		return &ResourceError{
			code: MissingRequiredAttributeErrCode,
			msg:  msg,
		}
	}
	InvalidRequestError = func(msg string) *ResourceError {
		return &ResourceError{
			code: InvalidRequestErrCode,
			msg:  msg,
		}
	}
)

type ResourceError struct {
	code     ErrorCode
	resource string
	ns       string
	name     string
	msg      string
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
	var resourcePath = ""
	if err.ns != "" || err.resource != "" {
		resourcePath = err.ns + " / " + err.resource
	}
	if err.name != "" {
		resourcePath += "( " + err.name + " )"
	}
	return fmt.Sprintf("%s [%v] %s", resourcePath, err.code, err.msg)
}
