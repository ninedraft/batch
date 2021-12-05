package batcher

func newBuffer[T any](size uint) *buffer[T] {
	return &buffer[T]{
		Items: make([]T, 0, size),
	}
}

type buffer[T any] struct {
	Items []T
}

func (b *buffer[T]) Reset() {
	if len(b.Items) == 0 {
		return
	}
	var empty T
	b.Items[0] = empty
	for i := 1; i < len(b.Items); i *= 2 {
		copy(b.Items[i:], b.Items[:i])
	}

	b.Items = b.Items[:0]
}

func (b *buffer[T]) Len() int {
	return len(b.Items)
}

func (b *buffer[T]) Push(v T) {
	b.Items = append(b.Items, v)
}
