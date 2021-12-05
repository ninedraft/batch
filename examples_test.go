package batch_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ninedraft/batch"
)

func ExampleBatcher() {
	var ctx = context.Background()

	var fn =  func(_ context.Context, values []string) error {
		fmt.Println(values)
		return nil
	}
	var b = batch.New(4,fn)

	for i := 0; i < 4; i++ {
		if err := b.Push(ctx, strconv.Itoa(i)); err != nil {
			panic(err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		panic(err)
	}
	// Output: [0 1 2 3]
}


func ExampleBatcher_Reset() {
	var ctx = context.Background()

	var fn =  func(_ context.Context, values []string) error {
		fmt.Println(values)
		return nil
	}
	var b = batch.New(4,fn)

	if err := b.Push(ctx, "bad balue"); err != nil {
		panic(err)
	}
	b.Reset(fn)
	if err := b.Push(ctx, "good value"); err != nil {
		panic(err)
	}
	if err := b.Flush(ctx); err != nil {
		panic(err)
	}
	// Output: [good value]
}

