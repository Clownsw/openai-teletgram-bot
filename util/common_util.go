package util

import "io"

// CommonClose close implement io.Closer object
func CommonClose[T io.Closer](t T) {
	_ = t.Close()
}
