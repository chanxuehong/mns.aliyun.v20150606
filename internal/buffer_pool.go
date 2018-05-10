package internal

import (
	"bytes"
	"sync"
)

var BufferPool = sync.Pool{
	New: newBuffer,
}

func newBuffer() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 16<<10))
}
