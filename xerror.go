package xerror

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pubgo/xerror/internal/color"
	"github.com/pubgo/xerror/xerror_abc"
)

type XErr = xerror_abc.XErr
type xerror struct {
	Err    error     `json:"cause,omitempty"`
	Msg    string    `json:"msg,omitempty"`
	Caller [2]string `json:"caller,omitempty"`
}

func (t *xerror) String() string            { return t.Stack() }
func (t *xerror) Debug(args ...interface{}) { p(handle(Wrap(t, args...)).stackString()) }
func (t *xerror) Unwrap() error             { return t.Err }
func (t *xerror) Cause() error              { return t.Err }
func (t *xerror) Wrap(args ...interface{}) error {
	return handle(t, func(err *xerror) { err.Msg = fmt.Sprint(args...) })
}

func (t *xerror) WrapF(msg string, args ...interface{}) error {
	return handle(t, func(err *xerror) { err.Msg = fmt.Sprintf(msg, args...) })
}

func (t *xerror) _p(buf *strings.Builder, xrr *xerror) {
	buf.WriteString("========================================================================================================================\n")
	if xrr.Err != nil {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Red.P("Err"), xrr.Err.Error()))
	}
	if strings.TrimSpace(xrr.Msg) != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("Msg"), xrr.Msg))
	}

	for i := range xrr.Caller {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Yellow.P("Caller"), xrr.Caller[i]))
	}

	if errs := trans(xrr.Err); errs != nil {
		for i := range errs {
			t._p(buf, errs[i])
		}
	}
}

func (t *xerror) stackString() string {
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

func (t *xerror) Is(err error) bool {
	if t == nil || t.Err == nil || err == nil {
		return false
	}

	switch err := err.(type) {
	case *xerrorBase:
		return err == t.Err
	case *xerror:
		return err == t || err.Err == t.Err
	case error:
		return t.Err == err
	default:
		return false
	}
}

func (t *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			type errors xerror
			_, _ = fmt.Fprintf(s, "%#v", (*errors)(t))
			return
		}

		if s.Flag('+') {
			_, _ = fmt.Fprint(s, t.Stack(true))
			return
		}

		_, _ = fmt.Fprint(s, t.Stack())
	case 's', 'q':
		_, _ = fmt.Fprint(s, t.Msg+": \n\t"+t.Error()+"\n\t"+t.Caller[0]+"\n\t"+t.Caller[1])
	default:
		_, _ = fmt.Fprint(s, t.Msg)
	}
}

func (t *xerror) Stack(indent ...bool) string {
	if t == nil || t.Err == nil || t.Err == ErrDone {
		return ""
	}

	if len(indent) > 0 {
		dt, err := json.MarshalIndent(t, "", "\t")
		if err != nil {
			log.Fatalln(err)
		}
		return string(dt)
	}
	dt, err := json.Marshal(t)
	if err != nil {
		log.Fatalln(err)
	}
	return string(dt)
}

func (t *xerror) As(target interface{}) bool {
	t1 := reflect.Indirect(reflect.ValueOf(target)).Interface()
	if err, ok := t1.(*xerror); ok {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

// Error
func (t *xerror) Error() string {
	if t == nil || isErrNil(t.Err) {
		return ""
	}

	return t.Err.Error()
}
