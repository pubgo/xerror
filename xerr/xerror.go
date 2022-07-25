package xerr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/settings"
)

const CallStackDepth = 2

func New(format string, a ...interface{}) XErr {
	x := &XError{}
	x.Msg = fmt.Sprintf(format, a...)
	x.Caller = []string{utils.CallerWithDepth(CallStackDepth + 1)}
	return x
}

type XError struct {
	Err    error    `json:"cause,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Detail string   `json:"detail,omitempty"`
	Caller []string `json:"caller,omitempty"`
}

func (t *XError) xErr()          {}
func (t *XError) String() string { return t.Stack() }
func (t *XError) DebugPrint() {
	if !settings.Debug {
		return
	}

	p(t.debugString())
	debug.PrintStack()
}

func (t *XError) Unwrap() error { return t.Err }
func (t *XError) Cause() error  { return t.Err }
func (t *XError) Wrap(args ...interface{}) XErr {
	return WrapXErr(t, func(err *XError) { err.Detail = fmt.Sprint(args...) })
}

func (t *XError) WrapF(msg string, args ...interface{}) XErr {
	return WrapXErr(t, func(err *XError) { err.Detail = fmt.Sprintf(msg, args...) })
}

func (t *XError) _p(buf *strings.Builder, xrr *XError) {
	if xrr == nil {
		return
	}

	buf.WriteString("========================================================================================================================\n")
	if xrr.Msg != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("Msg"), xrr.Msg))
	}

	if xrr.Detail != "" {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("Detail"), xrr.Detail))
	}

	for i := range xrr.Caller {
		if strings.Contains(xrr.Caller[i], "/src/runtime/") {
			continue
		}
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Yellow.P("Caller"), xrr.Caller[i]))
	}

	t._p(buf, trans(xrr.Err))
}

func (t *XError) debugString() string {
	if t == nil || t.Err == nil {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString("\n")
	t._p(buf, t)
	buf.WriteString("========================================================================================================================\n\n")
	return buf.String()
}

func (t *XError) Is(err error) bool {
	if t == nil || t.Err == nil || err == nil {
		return false
	}

	switch _err := err.(type) {
	case *XError:
		return _err == t || _err.Err == t.Err
	case error:
		return t.Err == _err
	default:
		return false
	}
}

func (t *XError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			type errors XError
			_, _ = fmt.Fprintf(s, "%#v", (*errors)(t))
			return
		}

		if s.Flag('+') {
			_, _ = fmt.Fprint(s, t.Stack())
			return
		}

		_, _ = fmt.Fprint(s, t.Stack())
	case 's', 'q':
		_, _ = fmt.Fprint(s, t.Msg+": \n\t"+t.Error()+"\n\t"+t.Caller[0]+"\n\t"+t.Caller[1])
	default:
		_, _ = fmt.Fprint(s, t.Msg)
	}
}

func (t *XError) Stack() string {
	if t == nil || t.Err == nil {
		return ""
	}

	dt, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(dt)
}

func (t *XError) As(target interface{}) bool {
	if t == nil || target == nil {
		return false
	}

	var v = reflect.ValueOf(target)
	t1 := reflect.Indirect(v).Interface()
	if err, ok := t1.(*XError); ok {
		v.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

// Error
func (t *XError) Error() string {
	if t == nil || isErrNil(t.Err) {
		return ""
	}

	return t.Err.Error()
}
