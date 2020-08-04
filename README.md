# xerror

go error 简单实现


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
package xerror_http

import (
	"github.com/pubgo/xerror"
	"net/http"
)

var (
	ErrHttp                = xerror.New("http error", "http错误")
	ErrBadRequest          = ErrHttp.New("400", http.StatusText(400))
	ErrUnauthorized        = ErrHttp.New("401", http.StatusText(401))
	ErrForbidden           = ErrHttp.New("403", http.StatusText(403))
	ErrNotFound            = ErrHttp.New("404", http.StatusText(404))
	ErrMethodNotAllowed    = ErrHttp.New("405", http.StatusText(405))
	ErrTimeout             = ErrHttp.New("408", http.StatusText(408))
	ErrConflict            = ErrHttp.New("409", http.StatusText(409))
	ErrInternalServerError = ErrHttp.New("500", http.StatusText(500))
)
```

```go
package xerror_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/xerror_core"
	"github.com/pubgo/xerror/xerror_http"
)

func check(b bool) {
	if !b {
		log.Fatalln("")
	}
}

func panic1(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(xerror_http.ErrBadRequest, "panic1 %+v", a)
	return
}

func panic2(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(panic1(a...), "panic2 %+v", a)
	return
}

func panicWrap(a ...interface{}) (err error) {
	return xerror.WrapF(panic2(a...), "panicWrap %+v", a)
}

func TestStack(t *testing.T) {
	defer xerror.Resp(func(err xerror.XErr) {
		fmt.Println(err.Stack(true))
	})
	xerror.Panic(panicWrap(1, 2, 4, 5))
}

func TestAs(t *testing.T) {
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrHttp) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrBadRequest) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrNotFound) == false)
}

func TestExit(t *testing.T) {
	xerror_core.PrintStack = false
	xerror.Exit(panicWrap(1, 2, 4, 5))
}

func TestTry(t *testing.T) {
	fmt.Println(xerror.Try(func() {
		panic("hello")
	}))
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(xerror_http.ErrBadRequest)
			return
		}()
	}
}

func BenchmarkPanicWithoutCaller(b *testing.B) {
	xerror_core.IsCaller = false
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(xerror_http.ErrBadRequest)
			return
		}()
	}
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