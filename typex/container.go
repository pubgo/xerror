package typex

type M[T any] map[string]T
type D[T any] []*kv[T]

// Map creates a map from the elements of the D.
func (d D[T]) Map() M[T] {
	m := make(M[T], len(d))
	for _, e := range d {
		m[e.K] = e.V
	}
	return m
}

func (d *D[T]) Append(kv ...*kv[T]) {
	*d = append(*d, kv...)
}

func KvOf[T any](k string, v T) *kv[T] {
	return &kv[T]{K: k, V: v}
}

// kv represents a BSON element for a D. It is usually used inside a D.
type kv[T any] struct {
	K string `json:"k"`
	V T      `json:"v"`
}

func (kv kv[T]) Map() M[T] { return M[T]{kv.K: kv.V} }
