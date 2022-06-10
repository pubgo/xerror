package main

import (
	"github.com/pubgo/xerror"
)

// 单个pkg的error处理

var err1 = &xerror.Err{Msg: "业务错误处理", Detail: "详细信息"}

func Hello() {
	defer xerror.RecoverAndRaise(func(err xerror.XErr) xerror.XErr {
		return err.Wrap("Hello wrap")
	})

	var err2 = xerror.WrapF(err1, "处理 wrap")
	xerror.Panic(err2, "处理 panic")
	return
}

func CallHello() (gErr error) {
	defer xerror.Recovery(func(err xerror.XErr) {
		gErr = err.WrapF("CallHello wrap")
	})

	Hello()

	return
}

func main() {
	defer xerror.RecoverAndExit()

	xerror.Panic(CallHello())
}
