package xerror_core

const Delimiter = "||"

var IsCaller bool
var CallDepth = 3
var PrintStack bool

func init() {
	IsCaller = true
	PrintStack = true
}
