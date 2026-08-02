package richerror

type Kind uint8

const (
	KindUnknown Kind = iota + 1

	KindNotFound

	KindValidation

	KindUnauthorized

	KindForbidden

	KindConflict

	KindUnexpected
)
