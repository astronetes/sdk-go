package errors

import (
	"errors"
	"fmt"
)

const errorsSite = "https://labs.astronetes.com/docs/errors"

type AstronetesError struct {
	code ErrorCode
	msg  string
	err  error
	meta map[string]any
}

func (err *AstronetesError) Msg() string {
	return err.msg
}

func (err *AstronetesError) DocRef() string {
	return fmt.Sprintf("%s/%s", errorsSite, err.code)
}

func (err *AstronetesError) Code() ErrorCode {
	return err.code
}

func (err *AstronetesError) Set(key string, value any) {
	if err.meta == nil {
		err.meta = make(map[string]any)
	}
	err.meta[key] = value
}

func (err *AstronetesError) Is(code ErrorCode) bool {
	return err.code == code
}

func (err *AstronetesError) Unwrap() error {
	return err.err
}

func (err *AstronetesError) Dig() error {
	var resErr *ResourceError
	if errors.As(err.err, &resErr) {
		return resErr
	}
	var ctrlErr *ControllerError
	if errors.As(err.err, &ctrlErr) {
		return ctrlErr
	}
	return err.err
}
