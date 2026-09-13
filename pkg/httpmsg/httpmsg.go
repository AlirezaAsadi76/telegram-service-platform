package httpmsg

import (
	"errors"
	"net/http"
	"telegram-service-platform/pkg/richerror"
)

func CodeAndMessage(err error) (message string, code int) {
	var richError richerror.RichError
	switch {
	case errors.As(err, &richError):
		var er richerror.RichError
		errors.As(err, &er)
		msg := er.Message()
		code := mapKindToHttpStatusCode(er.Kind())
		if code >= 500 {
			msg = "internal server error"
		}
		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func mapKindToHttpStatusCode(kind richerror.Kind) int {

	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
