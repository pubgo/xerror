package xerror_core

// errMsgHandlers 允许使用方自定义对recovery结果的处理
var errMsgHandlers []func(v interface{}) error

func Register(fn func(v interface{}) error) {
	errMsgHandlers = append(errMsgHandlers, fn)
}

func Handlers() []func(v interface{}) error {
	return errMsgHandlers
}
