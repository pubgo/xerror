package syncx

import (
	"context"
	"time"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/logx"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/xerr"
)

func Async[T any](fn func() typex.Value[T]) chan typex.Value[T] { return GoChan[T](fn) }

// GoChan 通过chan的方式同步执行异步任务
func GoChan[T any](fn func() typex.Value[T]) chan typex.Value[T] {
	assert.If(fn == nil, "[fn] is nil")

	var ch = make(chan typex.Value[T])

	go func() {
		defer close(ch)
		defer recovery.Recovery(func(err xerr.XErr) {
			ch <- typex.Err[T](err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
		})
		ch <- fn()
	}()

	return ch
}

// GoSafe 安全并发处理
func GoSafe(fn func(), catch ...func(err xerr.XErr) xerr.XErr) {
	assert.If(fn == nil, "[fn] is nil")
	go funk.TryAndLog(fn, catch...)
}

// GoCtx 可取消并发处理
func GoCtx(fn func(ctx context.Context), cb ...func(err xerr.XErr)) context.CancelFunc {
	assert.If(fn == nil, "[fn] is nil")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer recovery.Recovery(func(err xerr.XErr) {
			if len(cb) != 0 {
				cb[0](err)
				return
			}

			logx.Error(err, err.Error(), "fn", utils.CallerWithFunc(fn))
		})

		fn(ctx)
	}()

	return cancel
}

// GoDelay 异步延迟处理
func GoDelay(fn func(), durations ...time.Duration) {
	assert.Assert(fn == nil, "[fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	assert.Assert(dur == 0, "[dur] should not be 0")

	go funk.TryAndLog(fn)

	time.Sleep(dur)
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func() error) (gErr error) {
	assert.Assert(dur <= 0, "[Timeout] [dur] should not be less than zero")
	assert.Assert(fn == nil, "[Timeout] [fn] is nil")

	defer recovery.Err(&gErr, func(err xerr.XErr) xerr.XErr {
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
