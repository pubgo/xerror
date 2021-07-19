package xerror

var (
	Err       = New(Name)
	ErrType   = Err.New("type not match")
	ErrAssert = Err.New("assert true")
)
