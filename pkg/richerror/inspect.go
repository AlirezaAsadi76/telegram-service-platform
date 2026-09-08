package richerror

import "errors"

func IsCode(err error, code Code) bool {
	if err == nil {
		return false
	}

	var richErr *RichError

	if !errors.As(err, &richErr) {
		return false
	}

	return richErr.Code() == code
}
