package richerror

import (
	"errors"
)

type Op string
type RichError struct {
	kind         Kind
	code         Code
	message      string
	operation    Op
	wrappedError error
	meta         map[string]any
}

func New(operation Op, err error) *RichError {

	return &RichError{

		operation: operation,

		wrappedError: err,
	}

}

func (e *RichError) Error() string {

	if e.message != "" {

		return e.message

	}

	if e.wrappedError != nil {

		return e.wrappedError.Error()

	}

	return ""

}

func (e *RichError) Unwrap() error {

	return e.wrappedError

}

func (e *RichError) Kind() Kind {

	if e.kind != KindUnknown {

		return e.kind

	}

	var richErr *RichError

	if errors.As(e.wrappedError, &richErr) {

		return richErr.Kind()

	}

	return KindUnknown

}

func (e *RichError) Code() Code {

	if e.code != "" {

		return e.code

	}

	if richErr, ok := errors.AsType[*RichError](e.wrappedError); ok {

		return richErr.Code()

	}

	return CodeUnknown

}

func (e *RichError) Message() string {

	if e.message != "" {

		return e.message

	}

	var richErr *RichError

	if errors.As(e.wrappedError, &richErr) {

		return richErr.Message()

	}

	return ""

}

func (e *RichError) Operation() string {

	return string(e.operation)

}

func (e *RichError) Meta() map[string]any {

	return e.meta

}

func (e *RichError) WithKind(kind Kind) *RichError {

	e.kind = kind

	return e

}

func (e *RichError) WithCode(code Code) *RichError {

	e.code = code

	return e

}

func (e *RichError) WithMessage(message string) *RichError {

	e.message = message

	return e

}

func (e *RichError) WithMeta(meta map[string]any) *RichError {

	e.meta = meta

	return e

}
