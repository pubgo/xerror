package funkonf

var Conf = new(struct {
	EnableCaller bool
	PrintStack   bool
	Delimiter    string
	Debug        bool
})

func init() {
	Conf.EnableCaller = true
	Conf.PrintStack = true
	Conf.Delimiter = "||"
}
