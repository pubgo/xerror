package errors

import (
	"errors"
	"fmt"

	"github.com/pubgo/funk/xerr"
)

type XError = xerr.XError
type XErr = xerr.XErr
type Err = xerr.Err

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func New(format string, a ...interface{}) xerr.XErr {
	return xerr.New(format, a...)
}

func Wrap(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprint(args...) })
}

func WrapFn(err error, fn func(err xerr.XErr) xerr.XErr) error {
	if err == nil {
		return nil
	}

	return fn(xerr.WrapXErr(err))
}

func WrapF(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprintf(msg, args...) })
}

func ParseXErr(err error, fns ...func(err *xerr.XError)) *xerr.XError {
	return xerr.WrapXErr(err, fns...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Cause(err error) error {
	for {
		err1 := errors.Unwrap(err)
		if err1 == nil {
			return err
		}

		err = err1
	}
}
