package extras

func Reference[T any](a T) *T {
	return &a
}
