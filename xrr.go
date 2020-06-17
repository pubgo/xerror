package xerror

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type xerror struct {
	error
	xrr    error
	Code1  string                 `json:"code,omitempty"`
	Err    string                 `json:"err,omitempty"`
	Msg    string                 `json:"msg,omitempty"`
	Caller string                 `json:"caller,omitempty"`
	Attach map[string]interface{} `json:"attached,omitempty"`
	Sub    *xerror                `json:"sub,omitempty"`
}

func (t *xerror) Attached(k string, v interface{}) {
	if t.Attach == nil {
		t.Attach = map[string]interface{}{k: v}
		return
	}
	t.Attach[k] = v
}

func (t *xerror) New(code, msg string) XErr {
	return &xerror{Code1: t.Code1 + ": " + code, Msg: msg}
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

// Format...
func (t *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, t.Detail())
	case 's':
		if t.xrr != nil {
			io.WriteString(s, t.Error())
		}
	case 'q':
		fmt.Fprintf(s, "%q", t.Error())
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
	var dt []byte

	if Debug {
		dt, _ = json.MarshalIndent(t, "", "\t")
	} else {
		dt, _ = json.Marshal(t)
	}
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
