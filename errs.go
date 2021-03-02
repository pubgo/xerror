package xerror

import "errors"

var (
	ErrType   = errors.New("[xerror] type not match")
	ErrAssert = errors.New("[xerror] assert true")
)
