package xerror_conf

type conf struct {
	EnableCaller bool
	CallDepth    int
	PrintStack   bool
	Delimiter    string
}

var Conf = conf{
	EnableCaller: true,
	CallDepth:    2,
	PrintStack:   true,
	Delimiter:    "||",
}
