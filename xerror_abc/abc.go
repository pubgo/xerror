package xerror_abc

type XErr interface {
	Error() string
	Stack(indent ...bool) string
	Debug(args ...interface{})
	Cause() error
	Is(err error) bool
	As(val interface{}) bool
	Wrap(args ...interface{}) error
	WrapF(msg string, args ...interface{}) error
}
