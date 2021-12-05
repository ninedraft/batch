package batch_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	. "github.com/ninedraft/batch"

	"github.com/stretchr/testify/require"
)

func TestBatcher_HappyPath(test *testing.T) {
	test.Logf("checks if happy path emits correct values")

	// setup
	var ctx = context.Background()
	var expected = []string{"1", "two", "crow", "Ж", "", "📚", "\x11", "eight", "nein", "10"}
	const size = 4
	if len(expected) < size {
		test.Fatalf("number of input values must be greater than batch size")
	}

	var collected []string
	var b = New(size, func(_ context.Context, items []string) error {
		collected = append(collected, items...)
		return nil
	})

	// run
	for _, item := range expected {
		if err := b.Push(ctx, item); err != nil {
			test.Fatalf("bather.Push: unexpected error: %v", err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		test.Fatalf("bather.Flush: unexpected error: %v", err)
	}

	// check
	require.Equalf(test, expected, collected, "batcher wrote values")
}

func TestBatcher_DefaultSize(test *testing.T) {
	test.Logf("checks if batcher behaves correctly with default size")

	// setup
	var ctx = context.Background()
	var expected = []string{"1", "two", "crow", "Ж", "", "📚", "\x11", "eight", "nein", "10"}
	const size = 0

	var collected []string
	var b = New(size, func(_ context.Context, items []string) error {
		collected = append(collected, items...)
		return nil
	})

	// run
	for _, item := range expected {
		if err := b.Push(ctx, item); err != nil {
			test.Fatalf("bather.Push: unexpected error: %v", err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		test.Fatalf("bather.Flush: unexpected error: %v", err)
	}

	// check
	require.Equalf(test, expected, collected, "batcher wrote values")
}

func TestBatcher_ShortInput(test *testing.T) {
	test.Logf("checks if short input emits correct values")

	// setup
	var ctx = context.Background()
	var expected = []string{"1", "two"}
	const size = 4
	if len(expected) >= size {
		test.Fatalf("number of input values must be less than batch size")
	}

	var collected []string
	var b = New(size, func(_ context.Context, items []string) error {
		collected = append(collected, items...)
		return nil
	})

	// run
	for _, item := range expected {
		if err := b.Push(ctx, item); err != nil {
			test.Fatalf("bather.Push: unexpected error: %v", err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		test.Fatalf("bather.Flush: unexpected error: %v", err)
	}

	// check
	require.Equalf(test, expected, collected, "batcher wrote values")
}

func TestBatcher_CallbackError(test *testing.T) {
	test.Logf("checks if an error from callback passed out of batcher")

	// setup
	var ctx = context.Background()
	var errExpected = errors.New("test error")
	const size = 4

	counter := 0
	var b = New(size, func(_ context.Context, items []string) error {
		counter++
		if counter > 1 {
			test.Fatalf("callback must not be called after resulting error")
		}
		return errExpected
	})

	for i := 0; i < 10; i++ {
		value := strconv.Itoa(i)
		var err = b.Push(ctx, value)
		if i >= size {
			require.ErrorIsf(test, err, errExpected, "batcher.Push: resulting error: %v")
		}
	}
	var err = b.Flush(ctx)
	require.ErrorIs(test, err, errExpected, "batcher.Flush: resulting error")
}

func TestBatcher_Reset(test *testing.T) {
	test.Logf("checks if resets clears the internal buffer")

	// setup
	var ctx = context.Background()
	var expected = []string{"1", "two", "crow", "Ж", "", "📚", "\x11", "eight", "nein", "10"}
	const size = 4
	var chunkA = expected[:size-1]

	var collected []string
	var fn = func(_ context.Context, values []string) error {
		collected = append(collected, values...)
		return nil
	}
	var b = New(size, fn)
	for _, value := range chunkA {
		var err = b.Push(ctx, value)
		require.NoErrorf(test, err, "batcher.Push: resulting error")
	}
	b.Reset(fn)
	require.Emptyf(test, collected, "no items must be emitted by batcher before reset")

	for _, value := range expected {
		var err = b.Push(ctx, value)
		require.NoErrorf(test, err, "batcher.Push: resulting error")
	}
	if err := b.Flush(ctx); err != nil {
		test.Fatalf("bather.Flush: unexpected error: %v", err)
	}
	require.Equalf(test, expected, collected, "batcher wrote values")
}

func TestBatcher_ResetEmpty(test *testing.T) {
	test.Logf("ensures that .Reset doesn't call the callback function")

	const size = 4
	fn := func(_ context.Context, values []string) error {
		test.Fatalf("callback must be never called")
		return nil
	}
	var b = New(size, fn)
	b.Reset(fn)
}

func TestBatcher_FlushEmpty(test *testing.T) {
	test.Logf("ensures that .Reset doesn't call the callback function")

	var ctx = context.Background()
	const size = 4
	fn := func(_ context.Context, values []string) error {
		test.Fatalf("callback must be never called")
		return nil
	}
	var b = New(size, fn)
	err := b.Flush(ctx)
	require.NoError(test, err, "batcher.Flush: resulting error")
}

func TestBatcher_Context(test *testing.T) {
	test.Logf("ensures context is passed through method to callback")

	type key struct{}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, key{}, "test value")

	var counter int
	var b = New(1, func(ctx context.Context, _ []string) error {
		counter++
		v := ctx.Value(key{})
		require.Equalf(test, "test value", v, "context value")
		return nil
	})

	errPush := b.Push(ctx, "1")
	require.NoError(test, errPush, "batcher.Push: resulting error")
	errFlush := b.Flush(ctx)
	require.NoError(test, errFlush, "batcher.Flush: resulting error")
	require.NotEmpty(test, counter, "callback must be called at least once")
}

func TestBatcher_Size(test *testing.T) {
	test.Logf("checks if happy path emits correct values")

	// setup
	var ctx = context.Background()
	var values = []string{"1", "two", "crow", "Ж", "", "📚", "\x11", "eight", "nein", "10"}
	const size = 4

	var b = New(size, func(_ context.Context, items []string) error {
		return nil
	})
	require.EqualValuesf(test, size, b.Size(), "batcher.Size")
	// run
	for _, value := range values {
		if err := b.Push(ctx, value); err != nil {
			test.Fatalf("bather.Push: unexpected error: %v", err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		test.Fatalf("bather.Flush: unexpected error: %v", err)
	}
	require.EqualValuesf(test, size, b.Size(), "batcher.Size")
}
