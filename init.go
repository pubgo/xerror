package xerror

import "github.com/pubgo/xerror/xerror_core"

var (
	// ErrDone done
	ErrDone        = New("DONE")
	ErrUnknownType = New("unknown type")
	ErrNotFuncType = New("not func type")
)

func isCaller() bool {
	return xerror_core.IsCaller
}
func callDepth() int {
	return xerror_core.CallDepth
}
