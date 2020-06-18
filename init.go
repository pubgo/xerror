package xerror

// func caller depth
const (
	callDepth = 3
)

var (
	// ErrDone done
	ErrDone        = New("DONE")
	ErrUnknownType = New("unknown type")
	ErrNotFuncType = New("not func type")
)
