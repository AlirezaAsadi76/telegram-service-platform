package httpmsg

import (
	"errors"
	"net/http"

	"telegram-service-platform/pkg/richerror"
)

func CodeAndMessage(err error) (string, int) {
	if err == nil {
		return "", http.StatusOK
	}

	var richErr *richerror.RichError

	if !errors.As(err, &richErr) {
		return "internal server error", http.StatusInternalServerError
	}

	code := mapKindToHTTPStatusCode(richErr.Kind())

	if code >= http.StatusInternalServerError {
		return "internal server error", code
	}

	message := richErr.Message()

	if message == "" {
		message = "request failed"
	}

	return message, code
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity

	case richerror.KindNotFound:
		return http.StatusNotFound

	case richerror.KindConflict:
		return http.StatusConflict

	case richerror.KindForbidden:
		return http.StatusForbidden

	case richerror.KindUnexpected:
		return http.StatusServiceUnavailable

	default:
		return http.StatusInternalServerError
	}
}
