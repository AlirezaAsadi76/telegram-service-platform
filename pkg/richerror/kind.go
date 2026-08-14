package richerror

import "errors"

type Kind uint8

const (
	KindUnknown Kind = iota + 1

	KindNotFound

	KindValidation

	KindUnauthorized

	KindForbidden

	KindConflict

	KindUnexpected

	KindInvalid

	KindQueryFailure

	KindScanFailure

	KindInternal

	KindDependency

	KindInfrastructure
	KindSerializationFailure
	KindIdempotencyFailure
	KindCreateFailed
	KindExternalAPI
	KindRedisNil
)

func IsKind(
	err error,
	kind Kind,
) bool {

	var rErr *RichError

	if !errors.As(err, &rErr) {
		return false
	}

	return rErr.Kind() == kind
}
