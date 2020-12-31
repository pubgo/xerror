package xerror

import "strings"

type combine []*xerror

func (t combine) String() string {
	var result []string
	for i := range t {
		result = append(result, "["+(t)[i].String()+"]")
	}
	return strings.Join(result, ", ")
}

func (t combine) Error() string {
	var result []string
	for i := range t {
		result = append(result, "["+(t)[i].Error()+"]")
	}
	return strings.Join(result, ", ")
}
