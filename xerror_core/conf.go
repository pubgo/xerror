package xerror_core

type conf struct {
	IsCaller   bool
	CallDepth  int
	PrintStack bool
	Delimiter  string
}

var Conf = conf{
	IsCaller:   true,
	CallDepth:  3,
	PrintStack: true,
	Delimiter:  "||",
}
