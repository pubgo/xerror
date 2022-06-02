package main

import (
	"errors"

	"github.com/pubgo/xerror"
)

// 单个pkg的error处理

var err1 = errors.New("业务错误处理")
var err2 = errors.New("其他错误")

func Hello(flag bool) (err error) {
	defer xerror.RecoverErr(&err)

	if flag {
		xerror.Panic(err1, "处理 业务错误处理 失败")
	}

	xerror.Panic(err2, "处理 其他错误 失败")

	return
}

func CallHello(flag bool) (gErr error) {
	defer xerror.Recovery(func(err xerror.XErr) {
		// 跳过err1
		if errors.Is(err, err1) {
			return
		}

		gErr = err
	})

	xerror.Panic(Hello(flag))

	return
}

func main() {
	defer xerror.RecoverAndExit()

	xerror.Panic(CallHello(true))
	xerror.Panic(CallHello(false))
}
