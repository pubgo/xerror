package xerror

import (
	"strings"
)

type combine []*xerror

func (t combine) String() string {
	var b strings.Builder
	defer b.Reset()
	for i := range t {
		b.WriteString(t[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

func (t combine) Error() string {
	var b strings.Builder
	defer b.Reset()
	for i := range t {
		b.WriteString(t[i].Error())
		b.WriteString("\n")
	}
	return b.String()
}
