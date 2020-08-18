package xerror

import "strings"

type xerrorCombine []*xerror

func (t *xerrorCombine) String() string {
	var result []string
	for i := range *t {
		result = append(result, "["+(*t)[i].String()+"]")
	}
	return strings.Join(result, ", ")
}

func (t *xerrorCombine) Error() string {
	var result []string
	for i := range *t {
		result = append(result, "["+(*t)[i].Error()+"]")
	}
	return strings.Join(result, ", ")
}
