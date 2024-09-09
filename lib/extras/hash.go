package extras

type Hashable interface {
	Hash() string
}

type Hash[T Hashable] struct {
	s string
	h *T
}

func NewHash[T Hashable](m T) Hash[T] {
	return Hash[T]{
		s: m.Hash(),
		h: &m,
	}
}

func (h Hash[T]) FromHash() T {
	return *h.h
}
