package xerror

func TryCatch(fn func(), catch ...func(err error)) {
	Assert(fn == nil, "[fn] should not be nil")

	if len(catch) > 0 && catch[0] != nil {
		defer Resp(func(err XErr) { catch[0](err) })
	}

	fn()
}

func TryWith(err *error, fn func()) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespErr(err)
	fn()

	return
}

func TryThrow(fn func(), args ...interface{}) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespRaise(func(err XErr) error { return err.Wrap(args...) })
	fn()

	return
}

func Try(fn func()) (err error) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespErr(&err)

	fn()
	return
}
