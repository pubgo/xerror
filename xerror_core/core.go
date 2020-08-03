package xerror_core

var IsCaller bool
var CallDepth = 3
var PrintStack bool

func init() {
	IsCaller = true
	PrintStack = true
}
