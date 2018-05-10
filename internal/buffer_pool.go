package internal

import (
	"bytes"

	"github.com/chanxuehong/pool"
)

var BufferPool = pool.NewBytesBufferPool(100, newBuffer)

func newBuffer() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 16<<10))
}
