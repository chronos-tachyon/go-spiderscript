package exprtree

import (
	"fmt"
	"sync"
)

type String struct {
	Buffer *Buffer
	Offset uint
	Length uint
}

type Buffer struct {
	mu     sync.RWMutex
	interp *Interp
	bytes  []byte
	id     BufferID
}

type LenCap struct {
	Len uint
	Cap uint
}

func (interp *Interp) NewBuffer() *Buffer {
	id := interp.allocateBuffer()
	buf := &Buffer{
		interp: interp,
		bytes:  nil,
		id:     id,
	}
	interp.registerBuffer(buf)
	return buf
}

func (buf *Buffer) ID() BufferID {
	return buf.id
}

func (buf *Buffer) Interp() *Interp {
	return buf.interp
}

func (buf *Buffer) Clone() *Buffer {
	dupe := buf.interp.NewBuffer()
	locked(&dupe.mu, func() {
		_ = buf.WithReadLock(func(bytes []byte) error {
			dupe.bytes = make([]byte, len(bytes), cap(bytes))
			copy(dupe.bytes, bytes)
			return nil
		})
	})
	return dupe
}

func (buf *Buffer) Reset(shrink bool) {
	locked(&buf.mu, func() {
		if shrink {
			buf.bytes = nil
			return
		}
		buf.bytes = buf.bytes[:0]
	})
}

func (buf *Buffer) LenCap() LenCap {
	var value LenCap
	_ = buf.WithReadLock(func(bytes []byte) error {
		value = LenCap{
			Len: buf.LenLocked(),
			Cap: buf.CapLocked(),
		}
		return nil
	})
	return value
}

func (buf *Buffer) Len() uint {
	return buf.LenCap().Len
}

func (buf *Buffer) Cap() uint {
	return buf.LenCap().Cap
}

func (buf *Buffer) Bytes() []byte {
	var out []byte
	_ = buf.WithReadLock(func(bytes []byte) error {
		out = make([]byte, len(bytes))
		copy(out, bytes)
		return nil
	})
	return out
}

func (buf *Buffer) String() string {
	var str string
	_ = buf.WithReadLock(func(bytes []byte) error {
		str = string(bytes)
		return nil
	})
	return str
}

func (buf *Buffer) Grow(min uint) {
	locked(&buf.mu, func() {
		buf.GrowLocked(min)
	})
}

func (buf *Buffer) Truncate(n uint) {
	locked(&buf.mu, func() {
		buf.TruncateLocked(n)
	})
}

func (buf *Buffer) AppendBytes(bytes []byte) {
	locked(&buf.mu, func() {
		buf.AppendBytesLocked(bytes)
	})
}

func (buf *Buffer) AppendString(str string) {
	locked(&buf.mu, func() {
		buf.AppendStringLocked(str)
	})
}

func (buf *Buffer) AppendBuffer(other *Buffer) {
	_ = other.WithReadLock(func(bytes []byte) error {
		buf.AppendBytes(bytes)
		return nil
	})
}

func (buf *Buffer) WithWriteLock(fn func(bytes []byte) error) error {
	var err error
	locked(&buf.mu, func() {
		if buf.bytes == nil {
			buf.bytes = make([]byte, 0, 32)
		}
		err = fn(buf.bytes)
	})
	return err
}

func (buf *Buffer) WithReadLock(fn func(bytes []byte) error) error {
	var err error
	locked(buf.mu.RLocker(), func() {
		err = fn(buf.bytes)
	})
	return err
}

func (buf *Buffer) BytesLocked() []byte {
	return buf.bytes
}

func (buf *Buffer) LenCapLocked() LenCap {
	return LenCap{
		Len: buf.LenLocked(),
		Cap: buf.CapLocked(),
	}
}

func (buf *Buffer) LenLocked() uint {
	return uint(len(buf.bytes))
}

func (buf *Buffer) CapLocked() uint {
	return uint(cap(buf.bytes))
}

func (buf *Buffer) GrowLocked(min uint) {
	const lowBits = (^uint(0) >> 1)

	oldCap := uint(cap(buf.bytes))
	if oldCap >= min {
		return
	}

	oldLen := uint(len(buf.bytes))
	newLen := oldLen
	newCap := oldCap
	if newCap < 32 {
		newCap = 32
	}
	for newCap < min && newCap <= lowBits {
		newCap <<= 1
	}
	for newCap < min {
		newCap += 4096
	}

	orig := buf.bytes
	dupe := make([]byte, newLen, newCap)
	copy(dupe, orig)
	buf.bytes = dupe
}

func (buf *Buffer) TruncateLocked(n uint) {
	oldLen := uint(len(buf.bytes))
	if n == oldLen {
		return
	}
	if n > oldLen {
		buf.GrowLocked(n)
	}
	if n < oldLen {
		for i := n; i < oldLen; i++ {
			buf.bytes[i] = 0
		}
	}
	buf.bytes = buf.bytes[:n]
}

func (buf *Buffer) AppendBytesLocked(bytes []byte) {
	length := uint(len(bytes))
	if length == 0 {
		return
	}

	oldLen := uint(len(buf.bytes))
	newLen := oldLen + length

	buf.GrowLocked(newLen)
	buf.bytes = buf.bytes[:newLen]
	copy(buf.bytes[oldLen:newLen], bytes)
}

func (buf *Buffer) AppendStringLocked(str string) {
	buf.AppendBytesLocked([]byte(str))
}

var _ fmt.Stringer = (*Buffer)(nil)
