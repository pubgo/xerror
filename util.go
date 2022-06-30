package funk

func Last[T any](args []T) (t T) {
	if len(args) == 0 {
		return
	}

	return args[len(args)-1]
}

func Ternary[T any](ok bool, a T, b T) T {
	if ok {
		return a
	}
	return b
}

func If(ok bool, fn func()) {
	if ok {
		fn()
	}
}
