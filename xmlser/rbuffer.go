// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlser

import (
	"encoding/xml"
	"io"
	"reflect"
	"strconv"
)

// A RBuffer reads values from a Boost binary serialization stream.
type RBuffer struct {
	r   io.Reader
	err error

	types registry

	tok xml.Token
	dec *xml.Decoder
}

// NewRBuffer returns a new read-only buffer that reads from r.
func NewRBuffer(r io.Reader) *RBuffer {
	return &RBuffer{
		types: newRegistry(),
		dec:   xml.NewDecoder(r),
	}
}

func (r *RBuffer) Err() error { return r.err }

func (r *RBuffer) next() {
	if r.err != nil {
		return
	}
	r.tok, r.err = r.dec.Token()
}

func (r *RBuffer) ReadHeader() Header {
	var hdr Header
	if r.dec == nil {
		r.err = ErrNotBoost
		return hdr
	}

	if r.err != nil {
		return hdr
	}

	for {
		r.next()
		switch tok := r.tok.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case magicStartElement:
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "signature":
						if attr.Value != magicHeader {
							r.err = ErrNotBoost
							return hdr
						}
					case "version":
						v := 0
						v, r.err = strconv.Atoi(attr.Value)
						if r.err != nil {
							r.err = ErrInvalidHeader
							return hdr
						}
						hdr.Version = uint16(v)
						return hdr
					}
				}
			}
		}
	}

	r.err = ErrInvalidHeader
	return hdr
}

func (r *RBuffer) ReadTypeDescr(typ reflect.Type) TypeDescr {
	if dtype, ok := r.types[typ]; ok {
		return dtype
	}

	var dtype TypeDescr
	dtype.UnmarshalBoostXML(r)
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
	if r.err != nil {
		return ""
	}
	var v = ""
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadBool() bool {
	if r.err != nil {
		return false
	}
	var v bool
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadU8() uint8 {
	if r.err != nil {
		return 0
	}
	var v uint8
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadU16() uint16 {
	if r.err != nil {
		return 0
	}
	var v uint16
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadU32() uint32 {
	if r.err != nil {
		return 0
	}
	var v uint32
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadU64() uint64 {
	if r.err != nil {
		return 0
	}
	var v uint64
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadI8() int8 {
	if r.err != nil {
		return 0
	}
	var v int8
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadI16() int16 {
	if r.err != nil {
		return 0
	}
	var v int16
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadI32() int32 {
	if r.err != nil {
		return 0
	}
	var v int32
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadI64() int64 {
	if r.err != nil {
		return 0
	}
	var v int64
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadF32() float32 {
	if r.err != nil {
		return 0
	}
	var v float32
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadF64() float64 {
	if r.err != nil {
		return 0
	}
	var v float64
	r.err = r.dec.Decode(&v)
	return v
}

func (r *RBuffer) ReadC64() complex64 {
	if r.err != nil {
		return 0
	}
	var v c64Type
	r.err = r.dec.Decode(&v)
	return complex(v.Real, v.Imag)
}

func (r *RBuffer) ReadC128() complex128 {
	if r.err != nil {
		return 0
	}
	var v c128Type
	r.err = r.dec.Decode(&v)
	return complex(v.Real, v.Imag)
}

var (
	_ io.Reader = (*RBuffer)(nil)
)
