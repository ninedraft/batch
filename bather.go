package batcher

import (
	"context"
)

// DefaultSize is a default batch size.
const DefaultSize uint = 64

// Batcher implements buffering (similar to bufio.Writer).
// Pushed values are accumulated in an internal buffer
// until Size number of values is reached.
// Then accumulated values are passed into the provided callback.
// If callback returns an error, then all subsequent calls of Push and Flush
// will return the same error until the Bather is reset with Reset.
type Batcher[T any] struct {
	size uint
	b    *buffer[T]
	err  error
	fn   Callback[T]
}

// Callback describes a value consuming function.
// Values slice must be used only while function execution
// and must not be shared between calls.
type Callback[T any] func(ctx context.Context, values []T) error

// New creates a new batcher with defined size and callback.
// If size is 0, then DefaultSize is used.
// If fn callback is nil, then Batcher may panic with future methods calls.
func New[T any](size uint, fn Callback[T]) *Batcher[T] {
	if size == 0 {
		size = DefaultSize
	}
	b := &Batcher[T]{
		size: size,
		b:    newBuffer[T](size),
		fn:   fn,
	}
	return b
}

// Push puts message into internal buffer.
// It will call the callback with provided context and filled buffer if len(buffer) > size.
// If callback is called, then Push returns resulting error.
// If callback is nil, then Push will panic.
func (b *Batcher[T]) Push(ctx context.Context, v T) error {
	if b.err != nil {
		return b.err
	}
	b.b.Push(v)
	if uint(b.b.Len()) < b.size {
		return nil
	}
	return b.Flush(ctx)
}

// Reset drops internal buffer and error.
func (b *Batcher[T]) Reset(fn Callback[T]) {
	b.b.Reset()
	b.err = nil
	b.fn = fn
}

// Flush passes internal buffer into the callback function and resets buffer.
// It returns resulting error, if got one.
func (b *Batcher[T]) Flush(ctx context.Context) error {
	if b.err != nil {
		return b.err
	}
	err := b.fn(ctx, b.b.Items)
	b.b.Reset()
	if err != nil {
		b.err = err
	}
	return b.err
}

// Size returns internal buffer size.
func (b *Batcher[T]) Size() uint {
	return b.size
}