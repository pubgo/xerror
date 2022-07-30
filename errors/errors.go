package errors

import (
	"errors"

	"github.com/pubgo/funk/xerr"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func New(format string, a ...interface{}) xerr.XErr {
	return xerr.New(format, a...)
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
