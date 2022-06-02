package xerror

const Name = "xerror"

type XErr interface {
	Error() string
	String() string
	DebugPrint()
	Unwrap() error
	Wrap(args ...interface{}) error
	WrapF(msg string, args ...interface{}) error
}
