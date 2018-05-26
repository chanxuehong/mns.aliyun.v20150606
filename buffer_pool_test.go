package mns

import (
	"bytes"
	"testing"
)

func TestDefaultBytesBufferPool(t *testing.T) {
	// new
	func() {
		pool := GetBytesBufferPool()

		buffer := pool.Get()
		defer pool.Put(buffer)

		have := buffer.String()
		want := ""
		if have != want {
			t.Errorf("have:%s, want:%s", have, want)
			return
		}
	}()

	// clean up
	{
		pool := GetBytesBufferPool()
		for i := 0; i < 10; i++ {
			pool.Get()
		}
	}

	// new twice
	func() {
		pool := GetBytesBufferPool()

		buffer := pool.Get()
		defer pool.Put(buffer)

		have := buffer.String()
		want := ""
		if have != want {
			t.Errorf("have:%s, want:%s", have, want)
			return
		}
		buffer.WriteString("buffer")

		buffer2 := pool.Get()
		defer pool.Put(buffer2)

		have = buffer2.String()
		want = ""
		if have != want {
			t.Errorf("have:%s, want:%s", have, want)
			return
		}
	}()

	// clean up
	{
		pool := GetBytesBufferPool()
		for i := 0; i < 10; i++ {
			pool.Get()
		}
	}

	// reuse
	func() {
		func() {
			pool := GetBytesBufferPool()

			buffer := pool.Get()
			defer pool.Put(buffer)

			have := buffer.String()
			want := ""
			if have != want {
				t.Errorf("have:%s, want:%s", have, want)
				return
			}
			buffer.WriteString("buffer")
		}()

		func() {
			pool := GetBytesBufferPool()

			buffer := pool.Get()
			defer pool.Put(buffer)

			have := buffer.String()
			want := "buffer"
			if have != want {
				t.Errorf("have:%s, want:%s", have, want)
				return
			}
		}()
	}()

	// clean up
	{
		pool := GetBytesBufferPool()
		for i := 0; i < 10; i++ {
			pool.Get()
		}
	}

	// reuse and put nil
	func() {
		func() {
			pool := GetBytesBufferPool()

			buffer := pool.Get()
			defer pool.Put(buffer)

			have := buffer.String()
			want := ""
			if have != want {
				t.Errorf("have:%s, want:%s", have, want)
				return
			}
			buffer.WriteString("buffer")
		}()

		pool := GetBytesBufferPool()
		for i := 0; i < 10; i++ {
			pool.Put(nil)
		}

		func() {
			pool := GetBytesBufferPool()

			buffer := pool.Get()
			defer pool.Put(buffer)

			have := buffer.String()
			want := "buffer"
			if have != want {
				t.Errorf("have:%s, want:%s", have, want)
				return
			}
		}()
	}()
}

func Test_SetBytesBufferPool_GetBytesBufferPool(t *testing.T) {
	defer SetBytesBufferPool(_defaultBytesBufferPool)

	// get default BytesBufferPool
	pool := GetBytesBufferPool()
	if _, ok := pool.(*bytesBufferPool); !ok {
		t.Error("want type *bytesBufferPool")
		return
	}

	// SetBytesBufferPool with nil
	pool = nil
	SetBytesBufferPool(pool)
	pool = GetBytesBufferPool()
	if _, ok := pool.(*bytesBufferPool); !ok {
		t.Error("want type *bytesBufferPool")
		return
	}

	// SetBytesBufferPool with non-nil
	SetBytesBufferPool(&testBytesBufferPool{})
	pool = GetBytesBufferPool()
	if _, ok := pool.(*testBytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}

	// SetBytesBufferPool with nil
	pool = nil
	SetBytesBufferPool(pool)
	pool = GetBytesBufferPool()
	if _, ok := pool.(*testBytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}
}

type testBytesBufferPool struct{}

func (*testBytesBufferPool) Get() *bytes.Buffer  { return nil }
func (*testBytesBufferPool) Put(x *bytes.Buffer) {}
