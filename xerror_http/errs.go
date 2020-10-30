package xerror_http

import (
	"net/http"

	"github.com/pubgo/xerror"
)

var (
	ErrHttp                = xerror.New("http error", "http错误")
	ErrBadRequest          = ErrHttp.New("400", http.StatusText(400))
	ErrUnauthorized        = ErrHttp.New("401", http.StatusText(401))
	ErrForbidden           = ErrHttp.New("403", http.StatusText(403))
	ErrNotFound            = ErrHttp.New("404", http.StatusText(404))
	ErrMethodNotAllowed    = ErrHttp.New("405", http.StatusText(405))
	ErrTimeout             = ErrHttp.New("408", http.StatusText(408))
	ErrConflict            = ErrHttp.New("409", http.StatusText(409))
	ErrInternalServerError = ErrHttp.New("500", http.StatusText(500))
)
