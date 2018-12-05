// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmlser provides types to read and write XML archives from the C++
// Boost Serialization library.
//
// Writing values to an output binary archive can be done like so:
//
//  enc := xmlser.NewEncoder(w)
//  err := enc.Encode("hello")
//
// And reading values from an input binary archive:
//
//  dec := xmlser.NewDecoder(r)
//  str := ""
//  err := dec.Decode(&str)
//
// For more informations, look at the examples for Encoder, Decoder and read/write Buffer.
package xmlser // import "github.com/go-boostio/boostio/xmlser"

//go:generate go run ./testdata/gen-xml-archive.go

import (
	"reflect"

	"github.com/go-boostio/boostio"
	"github.com/pkg/errors"
)

const (
	magicStartElement = "boost_serialization"
	magicHeader       = "serialization::archive"
)

var (
	ErrNotBoost         = errors.New("xmlser: not a Boost XML archive")
	ErrInvalidHeader    = errors.New("xmlser: invalid Boost XML archive header")
	ErrInvalidTypeDescr = errors.New("xmlser: invalid Boost XML archive type descriptor")
	ErrTypeNotSupported = errors.New("xmlser: type not supported")
	ErrInvalidArrayLen  = errors.New("xmlser: invalid array type")
)

var (
	zeroHdr Header
	bserHdr = Header{Version: boostio.Version}
)

// Unmarshaler is the interface implemented by types that can unmarshal a
// Boost XML description of themselves.
type Unmarshaler interface {
	UnmarshalBoostXML(r *RBuffer) error
}

// Marshaler is the interface implemented by types that can marshal themselves
// into a valid Boost serialization XML archive.
type Marshaler interface {
	MarshalBoostXML(w *WBuffer) error
}

// Header describes a boost XML archive.
type Header struct {
	Version uint16
}

func (hdr Header) MarshalBoostXML(w *WBuffer) error {
	if w.err != nil {
		return w.err
	}
	w.WriteU16("version", hdr.Version)
	return w.err
}

func (hdr *Header) UnmarshalBoostXML(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	hdr.Version = r.ReadU16()
	return r.err
}

// TypeDescr describes an on-disk boost XML archive type.
type TypeDescr struct {
	Version uint32
	ID      int64
	Level   int64
}

func (dt TypeDescr) MarshalBoostXML(w *WBuffer) error {
	if w.err != nil {
		return w.err
	}
	w.WriteI64("class_id", dt.ID)
	w.WriteI64("tracking_level", dt.Level)
	w.WriteU32("version", dt.Version)
	return w.err
}

func (dt *TypeDescr) UnmarshalBoostXML(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	dt.ID = r.ReadI64()
	dt.Level = r.ReadI64()
	dt.Version = r.ReadU32()
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
