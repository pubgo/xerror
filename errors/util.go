package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func WithErr(err *Error, err1 error) error {
	if err1 == nil {
		return nil
	}

	return err.Err(err1)
}

// BadRequest generates a 400 error.
func BadRequest(err *Error) *Error {
	return err.WithCode(Http2Code(400))
}

// Unauthorized generates a 401 error.
func Unauthorized(err *Error) *Error {
	return err.WithCode(Http2Code(401))
}

// Forbidden generates a 403 error.
func Forbidden(err *Error) *Error {
	return err.WithCode(Http2Code(403))
}

// NotFound generates a 404 error.
func NotFound(err *Error) *Error {
	return err.WithCode(Http2Code(404))
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(err *Error) *Error {
	return err.WithCode(Http2Code(405))
}

// Timeout generates a 408 error.
func Timeout(err *Error) *Error {
	return err.WithCode(Http2Code(408))
}

// Conflict generates a 409 error.
func Conflict(err *Error) *Error {
	return err.WithCode(Http2Code(409))
}

// InternalServer generates a 500 error.
func InternalServer(err *Error) *Error {
	return err.WithCode(Http2Code(500))
}

func Http2Code(code int32) codes.Code {
	switch code {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusPreconditionFailed:
		return codes.FailedPrecondition
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	}
	return codes.Unknown
}
