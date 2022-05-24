// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binser provides types to read and write binary archives from the C++
// Boost Serialization library.
//
// Writing values to an output binary archive can be done like so:
//
//	enc := binser.NewEncoder(w)
//	err := enc.Encode("hello")
//
// And reading values from an input binary archive:
//
//	dec := binser.NewDecoder(r)
//	str := ""
//	err := dec.Decode(&str)
//
// For more informations, look at the examples for Encoder, Decoder and read/write Buffer.
package binser // import "github.com/go-boostio/boostio/binser"

//go:generate go run ./testdata/gen-binary-archive.go

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/go-boostio/boostio"
)

const (
	magicHeader = "serialization::archive"
)

var (
	ErrNotBoost         = errors.New("binser: not a Boost binary archive")
	ErrInvalidHeader    = errors.New("binser: invalid Boost binary archive header")
	ErrInvalidTypeDescr = errors.New("binser: invalid Boost binary archive type descriptor")
	ErrTypeNotSupported = errors.New("binser: type not supported")
	ErrInvalidArrayLen  = errors.New("binser: invalid array type")
)

// Arch describes the size of on-disk pointers.
type Arch byte

const ptrSize = 32 << uintptr(^uintptr(0)>>63)

const (
	ArchHW = Arch(ptrSize) // ArchHW enables writing/reading "native-bits" archives (ie: using the architecture's size).
	Arch32 = Arch(32)      // Arch32 enables writing/reading 32-bits archives
	Arch64 = Arch(64)      // Arch64 enables writing/reading 64-bits archives
)

// NewEncoder creates a new encoder.
func (a Arch) NewEncoder(w io.Writer) *Encoder {
	enc := newEncoder(w, a)
	enc.Header = a.Header()
	return enc
}

func (a Arch) Header() Header {
	switch a {
	case 0:
		return ArchHW.Header()
	case Arch32:
		return bser32Hdr
	case Arch64:
		return bser64Hdr
	default:
		panic(fmt.Errorf("binser: invalid architecture (size %d)", int(a)))
	}
}

var (
	zeroHdr   Header
	bser64Hdr = Header{
		Version: boostio.Version,
		Flags: binary.LittleEndian.Uint64([]byte{
			0x4, 0x8, // size of int, long
			0x4, 0x8, // size of float, double
			0x1, 0x0, 0x0, 0x0, // little-endian
		}),
	}
	bser32Hdr = Header{
		Version: boostio.Version,
		Flags: binary.LittleEndian.Uint64([]byte{
			0x4, 0x4, // size of int, long
			0x4, 0x8, // size of float, double
			0x1, 0x0, 0x0, 0x0, // little-endian
		}),
	}
)

// Unmarshaler is the interface implemented by types that can unmarshal a binary
// Boost description of themselves.
type Unmarshaler interface {
	UnmarshalBoost(r *RBuffer) error
}

// Marshaler is the interface implemented by types that can marshal themselves
// into a valid binary Boost serialization archive.
type Marshaler interface {
	MarshalBoost(w *WBuffer) error
}

// Header describes a binary boost archive.
type Header struct {
	Version uint16
	Flags   uint64
}

func (hdr Header) MarshalBoost(w *WBuffer) error {
	if w.err != nil {
		return w.err
	}
	w.WriteU16(hdr.Version)
	w.WriteU64(hdr.Flags)
	return w.err
}

func (hdr *Header) UnmarshalBoost(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	hdr.Version = r.ReadU16()
	hdr.Flags = r.ReadU64()
	return r.err
}

// TypeDescr describes an on-disk binary boost archive type.
type TypeDescr struct {
	Version uint32
	Flags   uint8
}

func (dt TypeDescr) MarshalBoost(w *WBuffer) error {
	if w.err != nil {
		return w.err
	}
	w.WriteU32(dt.Version)
	w.WriteU8(dt.Flags)
	return w.err
}

func (dt *TypeDescr) UnmarshalBoost(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	dt.Version = r.ReadU32()
	dt.Flags = r.ReadU8()
	return r.err
}

type registry map[reflect.Type]TypeDescr

func newRegistry() registry {
	return registry(map[reflect.Type]TypeDescr{
		reflect.TypeOf(false):          TypeDescr{},
		reflect.TypeOf(uint8(0)):       TypeDescr{},
		reflect.TypeOf(uint16(0)):      TypeDescr{},
		reflect.TypeOf(uint32(0)):      TypeDescr{},
		reflect.TypeOf(uint64(0)):      TypeDescr{},
		reflect.TypeOf(int8(0)):        TypeDescr{},
		reflect.TypeOf(int16(0)):       TypeDescr{},
		reflect.TypeOf(int32(0)):       TypeDescr{},
		reflect.TypeOf(int64(0)):       TypeDescr{},
		reflect.TypeOf(float32(0.0)):   TypeDescr{},
		reflect.TypeOf(float64(0.0)):   TypeDescr{},
		reflect.TypeOf(complex64(0)):   TypeDescr{},
		reflect.TypeOf(complex128(0)):  TypeDescr{},
		reflect.TypeOf(""):             TypeDescr{},
		reflect.TypeOf([]bool{}):       TypeDescr{},
		reflect.TypeOf([]uint8{}):      TypeDescr{},
		reflect.TypeOf([]uint16{}):     TypeDescr{},
		reflect.TypeOf([]uint32{}):     TypeDescr{},
		reflect.TypeOf([]uint64{}):     TypeDescr{},
		reflect.TypeOf([]int8{}):       TypeDescr{},
		reflect.TypeOf([]int16{}):      TypeDescr{},
		reflect.TypeOf([]int32{}):      TypeDescr{},
		reflect.TypeOf([]int64{}):      TypeDescr{},
		reflect.TypeOf([]float32{}):    TypeDescr{},
		reflect.TypeOf([]float64{}):    TypeDescr{},
		reflect.TypeOf([]complex64{}):  TypeDescr{},
		reflect.TypeOf([]complex128{}): TypeDescr{},
	})
}

func isCxxBoostBuiltin(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

var (
	_ Marshaler   = (*Header)(nil)
	_ Unmarshaler = (*Header)(nil)
	_ Marshaler   = (*TypeDescr)(nil)
	_ Unmarshaler = (*TypeDescr)(nil)
)
