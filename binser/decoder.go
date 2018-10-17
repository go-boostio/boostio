// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// A Decoder reads and decodes values from a Boost binary serialization stream.
type Decoder struct {
	r      *Reader
	Header Header
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder checks the stream has a correct Boost binary header.
func NewDecoder(r io.Reader) *Decoder {
	rr := NewReader(r)
	return &Decoder{r: rr, Header: rr.readHeader()}
}

// Decode reads the next value from its input and stores it in the
// value pointed to by ptr.
func (dec *Decoder) Decode(ptr interface{}) error {
	if dec.r.err != nil {
		return dec.r.err
	}
	//	if v, ok := ptr.(encoding.BinaryUnmarshaling); ok {
	//		_, dec.err =
	//	}
	rv := reflect.Indirect(reflect.ValueOf(ptr))
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(dec.r.ReadBool())
	case reflect.Int8:
		rv.SetInt(int64(dec.r.ReadI8()))
	case reflect.Int16:
		rv.SetInt(int64(dec.r.ReadI16()))
	case reflect.Int32:
		rv.SetInt(int64(dec.r.ReadI32()))
	case reflect.Int64:
		rv.SetInt(dec.r.ReadI64())
	case reflect.Uint8:
		rv.SetUint(uint64(dec.r.ReadU8()))
	case reflect.Uint16:
		rv.SetUint(uint64(dec.r.ReadU16()))
	case reflect.Uint32:
		rv.SetUint(uint64(dec.r.ReadU32()))
	case reflect.Uint64:
		rv.SetUint(dec.r.ReadU64())
	case reflect.Float32:
		rv.SetFloat(float64(dec.r.ReadF32()))
	case reflect.Float64:
		rv.SetFloat(dec.r.ReadF64())
	case reflect.String:
		rv.SetString(dec.r.ReadString())
	case reflect.Struct:
		/*vers*/ _ = dec.r.ReadU32()
		/*flag*/ _ = dec.r.ReadU8()
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			dec.Decode(rv.Field(i).Addr().Interface())
		}
	case reflect.Slice:
		n := dec.r.ReadU64()
		if len, n := rv.Len(), int(n); len < n {
			rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), n-len, n)))
		}
		for i := 0; i < int(n); i++ {
			e := rv.Index(i)
			dec.Decode(e.Addr().Interface()) // FIXME(sbinet): do not go through Decode each time
		}
	case reflect.Array:
		/*vers*/ _ = dec.r.ReadU32() // FIXME(sbinet): is it really the version?
		/*flag*/ _ = dec.r.ReadU8() // FIXME(sbinet): is it really some flag?
		n := int(dec.r.ReadU64())
		if n != rv.Type().Len() {
			return errors.Errorf("binser: invalid array type")
		}
		for i := 0; i < n; i++ {
			e := rv.Index(i)
			dec.Decode(e.Addr().Interface()) // FIXME(sbinet): do not go through Decode each time
		}
	case reflect.Map:
		/*vers*/ _ = dec.r.ReadU32() // FIXME(sbinet): is it really the version?
		/*flag*/ _ = dec.r.ReadU8() // FIXME(sbinet): is it really some flag?
		n := int(dec.r.ReadU64())
		_ = dec.r.ReadU64() // FIXME(sbinet): what is this ?
		_ = dec.r.ReadU8()  // FIXME(sbinet): ditto ?
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		for i := 0; i < n; i++ {
			k := reflect.New(kt)
			dec.Decode(k.Interface()) // FIXME(sbinet): do not go through Decode each time
			v := reflect.New(vt)
			dec.Decode(v.Interface()) // FIXME(sbinet): do not go through Decode each time
			rv.SetMapIndex(k.Elem(), v.Elem())
		}

	default:
		return errors.Errorf("boost: invalid type %T", ptr)
	}
	return dec.r.err
}
