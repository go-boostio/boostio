// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

// A RBuffer reads values from a Boost binary serialization stream.
type RBuffer struct {
	r   io.Reader
	err error
	buf []byte

	types registry
}

// NewRBuffer returns a new read-only buffer that reads from r.
func NewRBuffer(r io.Reader) *RBuffer {
	return &RBuffer{
		r:     r,
		buf:   make([]byte, 8),
		types: newRegistry(),
	}
}

func (r *RBuffer) Err() error { return r.err }

func (r *RBuffer) ReadHeader() Header {
	var hdr Header
	if r.r == nil {
		r.err = ErrNotBoost
		return hdr
	}

	if r.err != nil {
		return hdr
	}

	v := r.ReadString()
	if v != magicHeader {
		r.err = ErrNotBoost
		return hdr
	}

	hdr.UnmarshalBoost(r)
	if r.err != nil {
		r.err = ErrInvalidHeader
	}
	return hdr
}

func (r *RBuffer) ReadTypeDescr(typ reflect.Type) TypeDescr {
	if dtype, ok := r.types[typ]; ok {
		return dtype
	}

	var dtype TypeDescr
	dtype.UnmarshalBoost(r)
	switch r.err {
	case nil:
		r.types[typ] = dtype
	default:
		r.err = ErrInvalidTypeDescr
	}
	return dtype
}

func (r *RBuffer) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	var n int
	n, r.err = io.ReadFull(r.r, p)
	return n, r.err
}

func (r *RBuffer) ReadString() string {
	n := r.ReadU64()
	if n == 0 || r.err != nil {
		return ""
	}
	raw := make([]byte, n)
	_, r.err = io.ReadFull(r.r, raw)
	return string(raw)
}

func (r *RBuffer) ReadBool() bool {
	r.load(1)
	switch uint8(r.buf[0]) {
	case 0:
		return false
	default:
		return true
	}
}

func (r *RBuffer) ReadU8() uint8 {
	r.load(1)
	return uint8(r.buf[0])
}

func (r *RBuffer) ReadU16() uint16 {
	r.load(2)
	return binary.LittleEndian.Uint16(r.buf[:2])
}

func (r *RBuffer) ReadU32() uint32 {
	r.load(4)
	return binary.LittleEndian.Uint32(r.buf[:4])
}

func (r *RBuffer) ReadU64() uint64 {
	r.load(8)
	return binary.LittleEndian.Uint64(r.buf[:8])
}

func (r *RBuffer) ReadI8() int8 {
	r.load(1)
	return int8(r.buf[0])
}

func (r *RBuffer) ReadI16() int16 {
	r.load(2)
	return int16(binary.LittleEndian.Uint16(r.buf[:2]))
}

func (r *RBuffer) ReadI32() int32 {
	r.load(4)
	return int32(binary.LittleEndian.Uint32(r.buf[:4]))
}

func (r *RBuffer) ReadI64() int64 {
	r.load(8)
	return int64(binary.LittleEndian.Uint64(r.buf[:8]))
}

func (r *RBuffer) ReadF32() float32 {
	r.load(4)
	return math.Float32frombits(binary.LittleEndian.Uint32(r.buf[:4]))
}

func (r *RBuffer) ReadF64() float64 {
	r.load(8)
	return math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8]))
}

func (r *RBuffer) ReadC64() complex64 {
	r.load(8)
	v0 := math.Float32frombits(binary.LittleEndian.Uint32(r.buf[0:4]))
	v1 := math.Float32frombits(binary.LittleEndian.Uint32(r.buf[4:8]))
	return complex(v0, v1)
}

func (r *RBuffer) ReadC128() complex128 {
	r.load(8)
	v0 := math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8]))
	r.load(8)
	v1 := math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8]))
	return complex(v0, v1)
}

func (r *RBuffer) load(n int) {
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

var (
	_ io.Reader = (*RBuffer)(nil)
)
