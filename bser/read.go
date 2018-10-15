// Copyright 2018 The boostio Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bser provides types to read and write binary archives from the C++
// Boost Serialization library.
package bser // import "github.com/sbinet/boostio/bser"

//go:generate go run ./testdata/gen-binary-archive.go

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"

	"github.com/pkg/errors"
)

var (
	errNotBoost = errors.New("bser: not a Boost binary archive")
)

type Decoder struct {
	r   io.Reader
	err error
	buf []byte

	Header Header
}

type Header struct {
	Version uint16
	Flags   uint64
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	d := Decoder{r: r, buf: make([]byte, 8)}
	err := d.readHeader()
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (dec *Decoder) Err() error { return dec.err }

func (dec *Decoder) Decode(ptr interface{}) error {
	if dec.err != nil {
		return dec.err
	}
	//	if v, ok := ptr.(encoding.BinaryUnmarshaling); ok {
	//		_, dec.err =
	//	}
	rv := reflect.Indirect(reflect.ValueOf(ptr))
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(dec.ReadBool())
	case reflect.Int8:
		rv.SetInt(int64(dec.ReadI8()))
	case reflect.Int16:
		rv.SetInt(int64(dec.ReadI16()))
	case reflect.Int32:
		rv.SetInt(int64(dec.ReadI32()))
	case reflect.Int64:
		rv.SetInt(dec.ReadI64())
	case reflect.Uint8:
		rv.SetUint(uint64(dec.ReadU8()))
	case reflect.Uint16:
		rv.SetUint(uint64(dec.ReadU16()))
	case reflect.Uint32:
		rv.SetUint(uint64(dec.ReadU32()))
	case reflect.Uint64:
		rv.SetUint(dec.ReadU64())
	case reflect.Float32:
		rv.SetFloat(float64(dec.ReadF32()))
	case reflect.Float64:
		rv.SetFloat(dec.ReadF64())
	case reflect.String:
		rv.SetString(dec.ReadString())
	case reflect.Struct:
		/*vers*/ _ = dec.ReadU32()
		/*flag*/ _ = dec.ReadU8()
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			dec.err = dec.Decode(rv.Field(i).Addr().Interface())
			if dec.err != nil {
				return dec.err
			}
		}
	default:
		return fmt.Errorf("boost: invalid type %T", ptr)
	}
	return dec.err
}

func (d *Decoder) readHeader() error {
	v := d.ReadString()
	if d.err != nil {
		return d.err
	}
	if v != "serialization::archive" {
		return errNotBoost
	}
	d.Header.Version = d.ReadU16()
	d.Header.Flags = d.ReadU64()
	log.Printf("header: %#v", d.Header)
	return nil
}

func (dec *Decoder) Read(p []byte) (int, error) {
	if dec.err != nil {
		return 0, dec.err
	}
	var n int
	n, dec.err = io.ReadFull(dec.r, p)
	return n, dec.err
}

func (dec *Decoder) ReadString() string {
	n := dec.ReadU64()
	if n == 0 || dec.err != nil {
		return ""
	}
	raw := make([]byte, n)
	_, dec.err = io.ReadFull(dec.r, raw)
	return string(raw)
}

func (dec *Decoder) ReadBool() bool {
	dec.load(1)
	switch uint8(dec.buf[0]) {
	case 0:
		return false
	default:
		return true
	}
}

func (dec *Decoder) ReadU8() uint8 {
	dec.load(1)
	return uint8(dec.buf[0])
}

func (dec *Decoder) ReadU16() uint16 {
	dec.load(2)
	return binary.LittleEndian.Uint16(dec.buf[:2])
}

func (dec *Decoder) ReadU32() uint32 {
	dec.load(4)
	return binary.LittleEndian.Uint32(dec.buf[:4])
}

func (dec *Decoder) ReadU64() uint64 {
	dec.load(8)
	return binary.LittleEndian.Uint64(dec.buf[:8])
}

func (dec *Decoder) ReadI8() int8 {
	dec.load(1)
	return int8(dec.buf[0])
}

func (dec *Decoder) ReadI16() int16 {
	dec.load(2)
	return int16(binary.LittleEndian.Uint16(dec.buf[:2]))
}

func (dec *Decoder) ReadI32() int32 {
	dec.load(4)
	return int32(binary.LittleEndian.Uint32(dec.buf[:4]))
}

func (dec *Decoder) ReadI64() int64 {
	dec.load(8)
	return int64(binary.LittleEndian.Uint64(dec.buf[:8]))
}

func (dec *Decoder) ReadF32() float32 {
	dec.load(4)
	return math.Float32frombits(binary.LittleEndian.Uint32(dec.buf[:4]))
}

func (dec *Decoder) ReadF64() float64 {
	dec.load(8)
	return math.Float64frombits(binary.LittleEndian.Uint64(dec.buf[:8]))
}

func (dec *Decoder) load(n int) {
	if dec.err != nil {
		return
	}

	nn, err := io.ReadFull(dec.r, dec.buf[:n])
	if err != nil {
		dec.err = err
		return
	}

	if nn < n {
		dec.err = io.ErrUnexpectedEOF
	}
}
