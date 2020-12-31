package xerror

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/pubgo/xerror/xerror_util"
)

type xerrorOptions struct {
	depth int
	msg   string
}

type Option func(t *xerrorOptions)

func With(opts ...Option) XError {
	var opt xerrorOptions

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Next() XError {
	return With(WithCaller(1))
}

func WithCaller(depth int) Option {
	return func(t *xerrorOptions) {
		t.depth = depth
	}
}

func WithMsg(msg string, args ...interface{}) Option {
	return func(t *xerrorOptions) {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		t.msg = msg
	}
}

func (t xerrorOptions) Next() xerrorOptions {
	opts := t
	opts.depth += 1
	return opts
}

// Combine combine multiple errors
func Combine(errs ...error) error { return With(WithCaller(1)).Combine(errs...) }
func (t xerrorOptions) Combine(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var _errs combine
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		_errs = append(_errs, handle(errs[i], t))
	}

	if len(_errs) == 0 {
		return nil
	}

	return _errs
}

// Parse parse error to xerror
func Parse(err error) XErr { return With().Parse(err) }
func (t xerrorOptions) Parse(err error) XErr {
	if isErrNil(err) {
		return nil
	}

	return handle(err, t)
}

func Try(fn func()) (err error) { return With(WithCaller(1)).Try(fn) }
func (t xerrorOptions) Try(fn func()) (err error) {
	if fn == nil {
		return New("the [fn] parameters should not be nil")
	}

	defer Resp(func(err1 XErr) {
		err = WrapF(err1, xerror_util.CallerWithFunc(fn))
	})

	fn()
	return
}

func Done()                   { With(WithCaller(1)).Done() }
func (t xerrorOptions) Done() { panic(ErrDone) }

func Panic(err error) { With(WithCaller(1)).Panic(err) }
func (t xerrorOptions) Panic(err error) {
	if isErrNil(err) {
		return
	}
	panic(handle(err, t))
}

func PanicF(err error, msg string, args ...interface{}) { With(WithCaller(1)).PanicF(err, msg, args...) }
func (t xerrorOptions) PanicF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	WithMsg(msg, args...)(&t)
	panic(handle(err, t))
}

func Wrap(err error) error { return With(WithCaller(1)).Wrap(err) }
func (t xerrorOptions) Wrap(err error) error {
	if isErrNil(err) {
		return nil
	}
	return handle(err, t)
}

func WrapF(err error, msg string, args ...interface{}) error {
	return With(WithCaller(1)).WrapF(err, msg, args...)
}
func (t xerrorOptions) WrapF(err error, msg string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	WithMsg(msg, args...)(&t)
	return handle(err, t)
}

// PanicErr
func PanicErr(d1 interface{}, err error) interface{} { return With(WithCaller(1)).PanicErr(d1, err) }
func (t xerrorOptions) PanicErr(d1 interface{}, err error) interface{} {
	if isErrNil(err) {
		return d1
	}
	panic(handle(err, t))
}

func PanicBytes(d1 []byte, err error) []byte { return With(WithCaller(1)).PanicBytes(d1, err) }
func (t xerrorOptions) PanicBytes(d1 []byte, err error) []byte {
	if isErrNil(err) {
		return d1
	}
	panic(handle(err, t))
}

func PanicStr(d1 string, err error) string { return With(WithCaller(1)).PanicStr(d1, err) }
func (t xerrorOptions) PanicStr(d1 string, err error) string {
	if isErrNil(err) {
		return d1
	}
	panic(handle(err, t))
}

func PanicFile(d1 *os.File, err error) *os.File { return With(WithCaller(1)).PanicFile(d1, err) }
func (t xerrorOptions) PanicFile(d1 *os.File, err error) *os.File {
	if isErrNil(err) {
		return d1
	}
	panic(handle(err, t))
}

func PanicResponse(d1 *http.Response, err error) *http.Response {
	return With(WithCaller(1)).PanicResponse(d1, err)
}
func (t xerrorOptions) PanicResponse(d1 *http.Response, err error) *http.Response {
	if isErrNil(err) {
		return d1
	}
	panic(handle(err, t))
}

func Exit(err error) { With(WithCaller(1)).Exit(err) }
func (t xerrorOptions) Exit(err error) {
	if isErrNil(err) {
		return
	}

	p(handle(err, t).p())
	printStack()
	os.Exit(1)
}

// ExitF
func ExitF(err error, msg string, args ...interface{}) { Next().ExitF(err, msg, args...) }
func (t xerrorOptions) ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	WithMsg(msg, args...)(&t)
	p(handle(err, t).p())
	printStack()
	os.Exit(1)
}

// ExitErr
func ExitErr(dat interface{}, err error) interface{} { return With(WithCaller(1)).ExitErr(dat, err) }
func (t xerrorOptions) ExitErr(dat interface{}, err error) interface{} {
	if isErrNil(err) {
		return dat
	}

	p(handle(err, t).p())
	printStack()
	os.Exit(1)
	return nil
}

// FamilyAs Assert if *err belongs to *target's family
func FamilyAs(err error, target interface{}) bool { return With(WithCaller(1)).FamilyAs(err, target) }
func (t xerrorOptions) FamilyAs(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	for err != nil {
		if x, ok := err.(interface{ FamilyAs(interface{}) bool }); ok && x.FamilyAs(target) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}
