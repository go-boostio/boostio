// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binser provides types to read and write binary archives from the C++
// Boost Serialization library.
package binser // import "github.com/go-boostio/boostio/binser"

//go:generate go run ./testdata/gen-binary-archive.go

import (
	"github.com/pkg/errors"
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

var (
	zeroHdr Header
	bserHdr = Header{Version: 0x11, Flags: 0}
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

var (
	_ Marshaler   = (*Header)(nil)
	_ Unmarshaler = (*Header)(nil)
	_ Marshaler   = (*TypeDescr)(nil)
	_ Unmarshaler = (*TypeDescr)(nil)
)
