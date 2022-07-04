package xerr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/pubgo/funk/funkonf"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/internal/utils"
)

const CallStackDepth = 2

func New(format string, a ...interface{}) XErr {
	x := &Xerror{}
	x.Msg = fmt.Sprintf(format, a...)
	x.Caller = []string{utils.CallerWithDepth(CallStackDepth + 1)}
	return x
}

type Xerror struct {
	Err    error    `json:"cause,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Detail string   `json:"detail,omitempty"`
	Caller []string `json:"caller,omitempty"`
}

func (t *Xerror) xErr()          {}
func (t *Xerror) String() string { return t.Stack() }
func (t *Xerror) DebugPrint() {
	if !funkonf.Conf.Debug {
		return
	}

	p(t.debugString())
	debug.PrintStack()
}

func (t *Xerror) Unwrap() error { return t.Err }
func (t *Xerror) Cause() error  { return t.Err }
func (t *Xerror) Wrap(args ...interface{}) XErr {
	return WrapXErr(t, func(err *Xerror) { err.Detail = fmt.Sprint(args...) })
}

func (t *Xerror) WrapF(msg string, args ...interface{}) XErr {
	return WrapXErr(t, func(err *Xerror) { err.Detail = fmt.Sprintf(msg, args...) })
}

func (t *Xerror) _p(buf *strings.Builder, xrr *Xerror) {
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

func (t *Xerror) debugString() string {
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

func (t *Xerror) Is(err error) bool {
	if t == nil || t.Err == nil || err == nil {
		return false
	}

	switch _err := err.(type) {
	case *Xerror:
		return _err == t || _err.Err == t.Err
	case error:
		return t.Err == _err
	default:
		return false
	}
}

func (t *Xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			type errors Xerror
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

func (t *Xerror) Stack() string {
	if t == nil || t.Err == nil {
		return ""
	}

	dt, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(dt)
}

func (t *Xerror) As(target interface{}) bool {
	if t == nil || target == nil {
		return false
	}

	var v = reflect.ValueOf(target)
	t1 := reflect.Indirect(v).Interface()
	if err, ok := t1.(*Xerror); ok {
		v.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

// Error
func (t *Xerror) Error() string {
	if t == nil || isErrNil(t.Err) {
		return ""
	}

	return t.Err.Error()
}
