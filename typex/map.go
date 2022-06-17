package typex

import "sync"

var globalMutex sync.Mutex

type Map[T any] struct {
	done bool
	data map[string]T
}

func (t *Map[T]) once() {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if t.done {
		return
	}

	t.done = true
	t.data = make(map[string]T)
}

func (t *Map[T]) check() {
	if t.done {
		return
	}

	t.once()
}

func (t *Map[T]) Has(key string) bool {
	t.check()

	_, ok := t.data[key]
	return ok
}

func (t *Map[T]) Map() map[string]T {
	t.check()

	return t.data
}

func (t *Map[T]) Get(key string) T {
	t.check()

	val, ok := t.data[key]
	if ok {
		return val
	}

	return Zero[T]()
}

func (t *Map[T]) Load(key string) (T, bool) {
	t.check()

	val, ok := t.data[key]
	return val, ok
}

func (t *Map[T]) Keys() []string {
	t.check()

	var keys = make([]string, 0, len(t.data))
	for k := range t.data {
		keys = append(keys, k)
	}
	return keys
}

func (t *Map[T]) Each(fn func(name string, val T)) {
	t.check()

	for k, v := range t.data {
		fn(k, v)
	}
}

func (t *Map[T]) Set(key string, val T) {
	t.check()

	t.data[key] = val
}

func (t *Map[T]) Del(key string) {
	t.check()

	delete(t.data, key)
}
