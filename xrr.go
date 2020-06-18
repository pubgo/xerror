package xerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type xerror struct {
	error
	xrr    error
	Code1  string  `json:"code,omitempty"`
	Err    string  `json:"err,omitempty"`
	Msg    string  `json:"msg,omitempty"`
	Caller string  `json:"caller,omitempty"`
	Sub    *xerror `json:"sub,omitempty"`
}

func (t *xerror) New(ms ...string) XErr {
	if len(ms) == 0 {
		logger.Fatalln("the parameter cannot be empty")
	}

	var msg, code string
	switch len(ms) {
	case 1:
		code = ms[0]
	case 2:
		code, msg = ms[0], ms[1]
	}

	code = t.Code1 + ": " + code
	xw := &xerrorWrap{xerror: new(xerror)}
	xw.Code1 = code
	xw.Msg = msg
	xw.xrr = errors.New(code)
	xw.Caller = callerWithDepth(callDepth)

	return xw
}

func (t *xerror) Code() string {
	return t.Code1
}

func (t *xerror) next(e *xerror) {
	if t.Sub == nil {
		t.Sub = e
		return
	}
	t.Sub.next(e)
}

func (t *xerror) Unwrap() error {
	if t == nil {
		return nil
	}

	return t.xrr
}

func (t *xerror) p() string {
	if t == nil || t.xrr == nil {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	xrr := t
	for xrr != nil {
		buf.WriteString("========================================================================================================================\n")
		if xrr.xrr != nil {
			buf.WriteString(fmt.Sprintf("   %s]: %s\n", colorize("Err", colorRed), xrr.xrr))
		}
		if xrr.Msg != "" {
			buf.WriteString(fmt.Sprintf("   %s]: %s\n", colorize("Msg", colorGreen), xrr.Msg))
		}
		if xrr.Code1 != "" {
			buf.WriteString(fmt.Sprintf("  %s]: %s\n", colorize("Code", colorGreen), xrr.Code1))
		}
		buf.WriteString(fmt.Sprintf("%s]: %s\n", colorize("Caller", colorYellow), xrr.Caller))
		xrr = xrr.Sub
	}
	buf.WriteString("========================================================================================================================\n\n")
	return buf.String()
}

// Format...
func (t *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = io.WriteString(s, t.Detail())
	case 's':
		if t.xrr != nil {
			_, _ = io.WriteString(s, t.Error())
		}
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", t.Error())
	}
}

func (t *xerror) Is(err error) bool {
	if t == nil || t.xrr == nil || err == nil {
		return false
	}

	switch err := err.(type) {
	case *xerror:
		return err == t || err.xrr == t.xrr || err.Code1 == t.Code1
	case error:
		return t.xrr == err
	default:
		return false
	}
}

func (t *xerror) As(err interface{}) bool {
	if t == nil || t.xrr == nil || err == nil {
		return false
	}

	switch e := err.(type) {
	case *xerror:
		return strings.HasPrefix(t.Code1, e.Code1)
	case error:
		return strings.HasPrefix(t.Code1, e.Error())
	case string:
		return strings.HasPrefix(t.Code1, e)
	default:
		return false
	}
}

func (t *xerror) Detail() string {
	if t == nil || t.xrr == nil || t.xrr == ErrDone {
		return ""
	}

	t.Err = t.xrr.Error()
	dt, _ := json.Marshal(t)
	return string(dt)
}

// Error
func (t *xerror) Error() string {
	if t == nil || t.xrr == nil || t.xrr == ErrDone {
		return ""
	}

	return t.xrr.Error()
}

func (t *xerror) Reset() {
	t.xrr = nil
	t.Code1 = ""
	t.Err = ""
	t.Msg = ""
	t.Caller = ""
	if t.Sub == nil {
		putXerror(t)
		return
	}

	sub := t.Sub
	t.Sub = nil
	putXerror(t)
	sub.Reset()
}
