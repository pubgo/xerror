package xerror

var (
	// ErrDone done
	ErrDone        = New("DONE")
	ErrUnknownType = New("unknown type")
	ErrNotFuncType = New("not func type")
)
