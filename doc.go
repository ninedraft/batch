// Package batch contains a generic batch buffer,
// which cann accumulate multiple items into one slice and pass it into user defined callback.
//
// Known issues:
//	+ goreport card does not support generics (yet);
//	+ gomarkdoc does not support generics (yet);
//	+ doc.go.dev does not support generics (yet);
package batch