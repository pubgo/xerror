package xerror_abc

import (
	"net/http"
	"os"
)

type XErr interface {
	Error() string
	Stack(indent ...bool) string
	String() string
	Print()
	Unwrap() error
	Cause() error
	Is(err error) bool
	Wrap(args ...interface{}) error
	WrapF(msg string, args ...interface{}) error
}

type XError interface {
	Combine(errs ...error) error
	Parse(err error) XErr
	Try(fn func()) (err error)
	Panic(err error, args ...interface{})
	Done()
	PanicF(err error, msg string, args ...interface{})
	Wrap(err error, args ...interface{}) error
	WrapF(err error, msg string, args ...interface{}) error
	PanicErr(d1 interface{}, err error) interface{}
	PanicBytes(d1 []byte, err error) []byte
	PanicStr(d1 string, err error) string
	PanicFile(d1 *os.File, err error) *os.File
	PanicResponse(d1 *http.Response, err error) *http.Response
	ExitErr(dat interface{}, err error) interface{}
	ExitF(err error, msg string, args ...interface{})
	Exit(err error, args ...interface{})
	FamilyAs(err error, target interface{}) bool
}
