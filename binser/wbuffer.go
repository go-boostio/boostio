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

type WBuffer struct {
	w   io.Writer
	err error
	buf []byte

	types Registry
}

func NewWBuffer(w io.Writer) *WBuffer {
	return &WBuffer{
		w:     w,
		buf:   make([]byte, 8),
		types: NewRegistry(),
	}
}

func (w *WBuffer) Err() error { return w.err }

func (w *WBuffer) WriteHeader(hdr Header) error {
	w.err = hdr.MarshalBoost(w)
	return w.err
}

func (w *WBuffer) WriteTypeDescr(rt reflect.Type) error {
	dt, ok := w.types[rt]
	if ok {
		return nil
	}
	dt = TypeDescr{Version: 0, Flags: 0}
	w.types[rt] = dt
	w.err = dt.MarshalBoost(w)
	return w.err
}

func (w *WBuffer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	var n int
	n, w.err = w.w.Write(p)
	return n, w.err
}

func (w *WBuffer) WriteString(v string) error {
	if w.err != nil {
		return w.err
	}
	w.WriteU64(uint64(len(v)))
	_, w.err = w.w.Write([]byte(v))
	return w.err
}

func (w *WBuffer) WriteBool(v bool) error {
	if w.err != nil {
		return w.err
	}
	switch v {
	case false:
		w.buf[0] = 0
	default:
		w.buf[0] = 1
	}
	w.write(1)
	return w.err
}

func (w *WBuffer) WriteU8(v uint8) error {
	if w.err != nil {
		return w.err
	}
	w.buf[0] = v
	w.write(1)
	return w.err
}

func (w *WBuffer) WriteU16(v uint16) error {
	if w.err != nil {
		return w.err
	}
	const n = 2
	binary.LittleEndian.PutUint16(w.buf[:n], v)
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteU32(v uint32) error {
	if w.err != nil {
		return w.err
	}
	const n = 4
	binary.LittleEndian.PutUint32(w.buf[:n], v)
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteU64(v uint64) error {
	if w.err != nil {
		return w.err
	}
	const n = 8
	binary.LittleEndian.PutUint64(w.buf[:n], v)
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteI8(v int8) error {
	if w.err != nil {
		return w.err
	}
	w.buf[0] = uint8(v)
	w.write(1)
	return w.err
}

func (w *WBuffer) WriteI16(v int16) error {
	if w.err != nil {
		return w.err
	}
	const n = 2
	binary.LittleEndian.PutUint16(w.buf[:n], uint16(v))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteI32(v int32) error {
	if w.err != nil {
		return w.err
	}
	const n = 4
	binary.LittleEndian.PutUint32(w.buf[:n], uint32(v))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteI64(v int64) error {
	if w.err != nil {
		return w.err
	}
	const n = 8
	binary.LittleEndian.PutUint64(w.buf[:n], uint64(v))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteF32(v float32) error {
	if w.err != nil {
		return w.err
	}
	const n = 4
	binary.LittleEndian.PutUint32(w.buf[:n], math.Float32bits(v))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteF64(v float64) error {
	if w.err != nil {
		return w.err
	}
	const n = 8
	binary.LittleEndian.PutUint64(w.buf[:n], math.Float64bits(v))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteC64(v complex64) error {
	if w.err != nil {
		return w.err
	}
	const n = 8
	binary.LittleEndian.PutUint32(w.buf[:4], math.Float32bits(real(v)))
	binary.LittleEndian.PutUint32(w.buf[4:], math.Float32bits(imag(v)))
	w.write(n)
	return w.err
}

func (w *WBuffer) WriteC128(v complex128) error {
	if w.err != nil {
		return w.err
	}
	const n = 8
	binary.LittleEndian.PutUint64(w.buf[:n], math.Float64bits(real(v)))
	w.write(n)
	binary.LittleEndian.PutUint64(w.buf[:n], math.Float64bits(imag(v)))
	w.write(n)
	return w.err
}

func (w *WBuffer) write(n int) error {
	if w.err != nil {
		return w.err
	}

	var nn int
	nn, w.err = w.w.Write(w.buf[:n])
	if w.err == nil && nn < n {
		w.err = io.ErrShortWrite
	}
	return w.err
}

var (
	_ io.Writer = (*WBuffer)(nil)
)
