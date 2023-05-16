package errors

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	MissingRequiredAttributeError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(MissingRequiredAttributeErrCode, msg, args...)
	}
	InvalidRequestError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(InvalidRequestErrCode, msg, args...)
	}
	InvalidFormatError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(InvalidFormatErrCode, msg, args...)
	}
	UnknownResourceError = func(msg string, args ...any) *ResourceError {
		return NewResourceError(UnknownErrCode, msg, args...)
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

func (err *ResourceError) WithResourceDetails(obj v1.Object) *ResourceError {
	err.ns = obj.GetNamespace()
	err.name = obj.GetName()
	return err
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

func (err *ResourceError) Wrap(er error) *ResourceError {
	err.err = er
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
