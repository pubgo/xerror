package xerr

import "fmt"

func Wrap(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return WrapXErr(err, func(err *XError) { err.Detail = fmt.Sprint(args...) })
}

func WrapFn(err error, fn func(err XErr) XErr) error {
	if err == nil {
		return nil
	}

	return fn(WrapXErr(err))
}

func WrapF(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return WrapXErr(err, func(err *XError) { err.Detail = fmt.Sprintf(msg, args...) })
}
