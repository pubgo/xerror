package main

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/xerr"
)

// 单个pkg的error处理

var err1 = &xerr.Err{Msg: "业务错误处理", Detail: "详细信息"}

func Hello() {
	defer recovery.Raise(func(err xerr.XErr) xerr.XErr {
		return err.Wrap("Hello wrap")
	})

	var err2 = xerr.WrapF(err1, "处理 wrap")
	assert.MustF(err2, "处理 panic")
	return
}

func CallHello() (gErr error) {
	defer recovery.Recovery(func(err xerr.XErr) {
		gErr = err.WrapF("CallHello wrap")
	})

	Hello()

	return
}

func main() {
	defer recovery.Exit()

	assert.Must(CallHello())
}
