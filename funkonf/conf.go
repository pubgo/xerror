package funkonf

type conf struct {
	EnableCaller bool
	Debug        bool
}

var Conf = &conf{
	EnableCaller: true,
	Debug:        true,
}
