package funk

type XErr interface {
	xErr()
	Error() string
	String() string
	DebugPrint()
	Unwrap() error
	Wrap(args ...interface{}) XErr
	WrapF(msg string, args ...interface{}) XErr
}
