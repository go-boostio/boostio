// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser

import (
	"encoding/binary"
	"io"
	"math"
)

// A Reader reads values from a Boost binary serialization stream.
type Reader struct {
	r   io.Reader
	err error
	buf []byte
}

// NewReader returns a new reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r, buf: make([]byte, 8)}
}

func (r *Reader) Err() error { return r.err }

func (r *Reader) readHeader() Header {
	var hdr Header
	if r.r == nil {
		r.err = ErrNotBoost
		return hdr
	}

	if r.err != nil {
		return hdr
	}

	v := r.ReadString()
	if v != "serialization::archive" {
		r.err = ErrNotBoost
		return hdr
	}
	hdr.Version = r.ReadU16()
	hdr.Flags = r.ReadU64()
	if r.err != nil {
		r.err = ErrInvalidHeader
	}
	return hdr
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	var n int
	n, r.err = io.ReadFull(r.r, p)
	return n, r.err
}

func (r *Reader) ReadString() string {
	n := r.ReadU64()
	if n == 0 || r.err != nil {
		return ""
	}
	raw := make([]byte, n)
	_, r.err = io.ReadFull(r.r, raw)
	return string(raw)
}

func (r *Reader) ReadBool() bool {
	r.load(1)
	switch uint8(r.buf[0]) {
	case 0:
		return false
	default:
		return true
	}
}

func (r *Reader) ReadU8() uint8 {
	r.load(1)
	return uint8(r.buf[0])
}

func (r *Reader) ReadU16() uint16 {
	r.load(2)
	return binary.LittleEndian.Uint16(r.buf[:2])
}

func (r *Reader) ReadU32() uint32 {
	r.load(4)
	return binary.LittleEndian.Uint32(r.buf[:4])
}

func (r *Reader) ReadU64() uint64 {
	r.load(8)
	return binary.LittleEndian.Uint64(r.buf[:8])
}

func (r *Reader) ReadI8() int8 {
	r.load(1)
	return int8(r.buf[0])
}

func (r *Reader) ReadI16() int16 {
	r.load(2)
	return int16(binary.LittleEndian.Uint16(r.buf[:2]))
}

func (r *Reader) ReadI32() int32 {
	r.load(4)
	return int32(binary.LittleEndian.Uint32(r.buf[:4]))
}

func (r *Reader) ReadI64() int64 {
	r.load(8)
	return int64(binary.LittleEndian.Uint64(r.buf[:8]))
}

func (r *Reader) ReadF32() float32 {
	r.load(4)
	return math.Float32frombits(binary.LittleEndian.Uint32(r.buf[:4]))
}

func (r *Reader) ReadF64() float64 {
	r.load(8)
	return math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8]))
}

func (r *Reader) load(n int) {
	if r.err != nil {
		return
	}

	nn, err := io.ReadFull(r.r, r.buf[:n])
	if err != nil {
		r.err = err
		return
	}

	if nn < n {
		r.err = io.ErrUnexpectedEOF
	}
}
