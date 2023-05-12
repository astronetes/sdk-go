package errors

import "fmt"

const errorsSite = "https://labs.astronetes.com/docs/errors"

type AstronetesError struct {
	code ErrorCode
	msg  string
	err  error
	meta map[string]any
}

func (err *AstronetesError) DocRef() string {
	return fmt.Sprintf("%s/%s", errorsSite, err.code)
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
