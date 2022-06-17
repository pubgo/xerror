package funkonf

var Conf = new(struct {
	EnableCaller bool
	CallDepth    int
	PrintStack   bool
	Delimiter    string
	Debug        bool
})

func init() {
	Conf.EnableCaller = true
	Conf.CallDepth = 2
	Conf.PrintStack = true
	Conf.Delimiter = "||"
}
