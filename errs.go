package xerror

import "errors"

var (
	// ErrDone done
	ErrDone   = errors.New("[xerror] done")
	ErrType   = errors.New("[xerror] type not match")
	ErrAssert = errors.New("[xerror] assert true")
)
