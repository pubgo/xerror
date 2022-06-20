package syncx

import (
	"context"
	"time"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/internal/utils"
	"k8s.io/klog/v2"
)

// GoChan 通过chan的方式同步执行异步任务
func GoChan[T any](fn func() *Value[T]) chan *Value[T] {
	funk.Assert(fn == nil, "[fn] is nil")

	var ch = make(chan *Value[T])

	go func() {
		defer close(ch)
		defer funk.Recovery(func(err funk.XErr) {
			ch <- Err[T](err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
		})

		if val := fn(); val == nil {
			ch <- new(Value[T])
		} else {
			ch <- val
		}
	}()

	return ch
}

// GoSafe 安全并发处理
func GoSafe(fn func(), catch ...func(err funk.XErr) funk.XErr) {
	funk.Assert(fn == nil, "[fn] is nil")
	go funk.TryAndLog(fn, catch...)
}

// GoCtx 可取消并发处理
func GoCtx(fn func(ctx context.Context), cb ...func(err funk.XErr)) context.CancelFunc {
	funk.Assert(fn == nil, "[fn] is nil")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer funk.Recovery(func(err funk.XErr) {
			if len(cb) != 0 {
				cb[0](err)
				return
			}

			klog.ErrorS(err, err.Error(), "fn", utils.CallerWithFunc(fn))
		})

		fn(ctx)
	}()

	return cancel
}

// GoDelay 异步延迟处理
func GoDelay(fn func(), durations ...time.Duration) {
	funk.Assert(fn == nil, "[fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	funk.Assert(dur == 0, "[dur] should not be 0")

	go funk.TryAndLog(fn)

	time.Sleep(dur)
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func()) (gErr error) {
	funk.Assert(dur <= 0, "[Timeout] [dur] should not be less than zero")
	funk.Assert(fn == nil, "[Timeout] [fn] is nil")

	defer funk.RecoverErr(&gErr, func(err funk.XErr) funk.XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	var done = make(chan struct{})

	go func() {
		defer close(done)
		gErr = funk.Try(fn)
	}()

	select {
	case <-time.After(dur):
		return context.DeadlineExceeded
	case <-done:
		return
	}
}
