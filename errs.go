package xerror

import "errors"

var (
	// ErrDone done
	ErrDone = errors.New("DONE")
	ErrType = errors.New("type error")
)
