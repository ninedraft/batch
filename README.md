<!-- Code generated by gomarkdoc. DO NOT EDIT -->

[![Go Report Card](https://goreportcard.com/badge/github.com/ninedraft/batch)](https://goreportcard.com/report/github.com/ninedraft/batch)
[![GitHub Actions](https://github.com/ninedraft/batch/workflows/Test/badge.svg)](https://github.com/ninedraft/batch/actions?query=workflow%3ATest+branch%3Amaster)
[![Go Reference](https://pkg.go.dev/badge/github.com/ninedraft/batch.svg)](https://pkg.go.dev/github.com/ninedraft/batch)

# batch

```go
import "github.com/ninedraft/batch"
```

Package batch contains a generic batch buffer\, which accumulates multiple items into one slice and pass it into user defined callback\.

### Known issues

```
- goreport card does not support generics (yet);
- gomarkdoc does not support generics (yet);
- doc.go.dev does not support generics (yet);
```

### Code quality politics

```
- no external non-test dependencies;
- code coverage >= 90% (~100% currently);
```

## Index

- [Constants](<#constants>)
- [type Batcher](<#type-batcher>)
  - [func New[T any](size uint, fn Callback[T]) *Batcher[T]](<#func-new>)
  - [func (b *Batcher[T]) Flush(ctx context.Context) error](<#func-badrecv-flush>)
  - [func (b *Batcher[T]) Push(ctx context.Context, v T) error](<#func-badrecv-push>)
  - [func (b *Batcher[T]) Reset(fn Callback[T])](<#func-badrecv-reset>)
  - [func (b *Batcher[T]) Size() uint](<#func-badrecv-size>)
- [type Callback](<#type-callback>)


## Constants

DefaultSize is a default batch size\.

```go
const DefaultSize uint = 64
```

## type [Batcher](<https://github.com/ninedraft/batch/blob/master/bather.go#L16-L21>)

Batcher implements buffering \(similar to bufio\.Writer\)\. Pushed values are accumulated in an internal buffer until Size number of values is reached\. Then accumulated values are passed into the provided callback\. If callback returns an error\, then all subsequent calls of Push and Flush will return the same error until the Bather is reset with Reset\.

```go
type Batcher[T any] struct {
    // contains filtered or unexported fields
}
```

<details><summary>Example</summary>
<p>

```go
package main

import (
	"context"
	"fmt"
	"github.com/ninedraft/batch"
	"strconv"
)

func main() {
	var ctx = context.Background()

	var fn = func(_ context.Context, values []string) error {
		fmt.Println(values)
		return nil
	}
	var b = batch.New(4, fn)

	for i := 0; i < 4; i++ {
		if err := b.Push(ctx, strconv.Itoa(i)); err != nil {
			panic(err)
		}
	}
	if err := b.Flush(ctx); err != nil {
		panic(err)
	}
}
```

#### Output

```
[0 1 2 3]
```

</p>
</details>

### func [New](<https://github.com/ninedraft/batch/blob/master/bather.go#L31>)

```go
func New[T any](size uint, fn Callback[T]) *Batcher[T]
```

New creates a new batcher with defined size and callback\. If size is 0\, then DefaultSize is used\. If fn callback is nil\, then Batcher may panic with future methods calls\.

### func \(\*BADRECV\) [Flush](<https://github.com/ninedraft/batch/blob/master/bather.go#L67>)

```go
func (b *Batcher[T]) Flush(ctx context.Context) error
```

Flush passes internal buffer into the callback function and resets buffer\. It returns resulting error\, if got one\.

### func \(\*BADRECV\) [Push](<https://github.com/ninedraft/batch/blob/master/bather.go#L47>)

```go
func (b *Batcher[T]) Push(ctx context.Context, v T) error
```

Push puts message into internal buffer\. It will call the callback with provided context and filled buffer if len\(buffer\) \> size\. If callback is called\, then Push returns resulting error\. If callback is nil\, then Push will panic\.

### func \(\*BADRECV\) [Reset](<https://github.com/ninedraft/batch/blob/master/bather.go#L59>)

```go
func (b *Batcher[T]) Reset(fn Callback[T])
```

Reset drops internal buffer and error\.

### func \(\*BADRECV\) [Size](<https://github.com/ninedraft/batch/blob/master/bather.go#L80>)

```go
func (b *Batcher[T]) Size() uint
```

Size returns internal buffer size\.

## type [Callback](<https://github.com/ninedraft/batch/blob/master/bather.go#L26>)

Callback describes a value consuming function\. Values slice must be used only while function execution and must not be shared between calls\.

```go
type Callback[T any] func(ctx context.Context, values []T) error
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
