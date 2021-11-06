# xerror

> go error 简单实现

1. 高效处理golang的error, 避免处理大量的 `err!=nil` 判断
2. 高效处理golang的recover, 让错误中包含丰富的堆栈信息
3. xerror实现标准As,Is,Unwrap接口, 可以和其他error库一起使用
4. 简单易用


## 性能分析
```sh
go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out ./...
go tool pprof -http=":8081" profile.out
go tool pprof -http=":8081" memprofile.out
```

```md
goos: darwin
goarch: amd64
pkg: github.com/pubgo/xerror
BenchmarkPanic-8                	 1353934	       883 ns/op	     128 B/op	       2 allocs/op
BenchmarkPanicWithOutCaller-8   	 3861938	       309 ns/op	      48 B/op	       1 allocs/op
BenchmarkNoPanic-8              	201641330	         5.87 ns/op	       0 B/op	       0 allocs/op
PASS
ok      github.com/pubgo/xerror 4.363s
```

## example
```go
package xerror_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/xerror"
)

func TestErr(t *testing.T) {
	fmt.Println(xerror.Wrap(xerror.ErrAssert))
}

func TestParseWith(t *testing.T) {
	var err = fmt.Errorf("hello error")
	xerror.ParseWith(err, func(err xerror.XErr) {
		fmt.Printf("%v\n", err)
	})
}

func TestRespTest(t *testing.T) {
	defer xerror.RespTest(t)
	TestPanic1(t)
}

func TestRespNext(t *testing.T) {
	defer xerror.RespExit("TestRespNext")
	TestPanic1(t)
}

func TestPanic1(t *testing.T) {
	//defer xerror.RespExit()
	defer xerror.RespRaise(func(err xerror.XErr) error {
		return xerror.WrapF(err, "test raise")
	})

	//xerror.Panic(xerror.New("ok"))
	xerror.Panic(fmt.Errorf("ss"))
}

func init1Next() (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(fmt.Errorf("test next"))
	return nil
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(nil)
			return
		}()
	}
}
```