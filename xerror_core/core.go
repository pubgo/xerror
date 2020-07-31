package xerror_core

var IsCaller bool
var CallDepth = 3

func init() {
	IsCaller = true
}
