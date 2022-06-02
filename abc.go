package xerror

const Name = "xerror"

type XErr interface {
	Error() string
	String() string
	DebugPrint()
	Unwrap() error
	Wrap(args ...interface{}) XErr
	WrapF(msg string, args ...interface{}) XErr
}
