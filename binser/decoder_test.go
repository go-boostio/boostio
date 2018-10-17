// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser_test

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/go-boostio/boostio/binser"
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := binser.NewDecoder(f)
	for _, tc := range typeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.New(reflect.TypeOf(tc.want)).Elem()
			if rv.Kind() == reflect.Map {
				rv.Set(reflect.MakeMap(rv.Type()))
			}
			err := dec.Decode(rv.Addr().Interface())
			if err != nil {
				t.Fatalf("could not read %q: %v", tc.name, err)
			}
			if got, want := rv.Interface(), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got=%#v (%T)\nwant=%#v (%T)", got, got, want, want)
			}
		})
	}
}

func TestInvalidArchive(t *testing.T) {
	for _, tc := range []struct {
		raw []byte
		err error
		val interface{}
	}{
		{
			raw: nil,
			err: binser.ErrNotBoost,
		},
		{
			raw: []byte("boost"),
			err: binser.ErrNotBoost,
		},
		{
			raw: []byte{5, 0, 0, 0, 0, 0, 0, 0, 'b', 'o', 'o', 's', 't'},
			err: binser.ErrNotBoost,
		},
		{
			raw: []byte{
				0x16, 0, 0, 0, 0, 0, 0, 0,
				's', 'e', 'r', 'i', 'a', 'l', 'i', 'z', 'a', 't', 'i', 'o', 'n',
				':', ':',
				'a', 'r', 'c', 'h', 'i', 'v', 'e',
			},
			err: binser.ErrInvalidHeader,
		},
		{
			raw: []byte{
				0x16, 0, 0, 0, 0, 0, 0, 0,
				's', 'e', 'r', 'i', 'a', 'l', 'i', 'z', 'a', 't', 'i', 'o', 'n',
				':', ':',
				'a', 'r', 'c', 'h', 'i', 'v', 'e',
				0,
			},
			err: binser.ErrInvalidHeader,
		},
		{
			raw: []byte{
				0x16, 0, 0, 0, 0, 0, 0, 0,
				's', 'e', 'r', 'i', 'a', 'l', 'i', 'z', 'a', 't', 'i', 'o', 'n',
				':', ':',
				'a', 'r', 'c', 'h', 'i', 'v', 'e',
				1, 0,
			},
			err: binser.ErrInvalidHeader,
		},
		{
			raw: []byte{
				0x16, 0, 0, 0, 0, 0, 0, 0,
				's', 'e', 'r', 'i', 'a', 'l', 'i', 'z', 'a', 't', 'i', 'o', 'n',
				':', ':',
				'a', 'r', 'c', 'h', 'i', 'v', 'e',
				1, 0,
				0,
			},
			err: binser.ErrInvalidHeader,
		},
		{
			raw: []byte{
				0x16, 0, 0, 0, 0, 0, 0, 0,
				's', 'e', 'r', 'i', 'a', 'l', 'i', 'z', 'a', 't', 'i', 'o', 'n',
				':', ':',
				'a', 'r', 'c', 'h', 'i', 'v', 'e',
				1, 0,
				1, 0, 0, 0, 0, 0, 0, 0, 1,
			},
			err: io.ErrUnexpectedEOF,
			val: new(uint16),
		},
	} {
		t.Run("", func(t *testing.T) {
			dec := binser.NewDecoder(bytes.NewReader(tc.raw))
			err := dec.Decode(tc.val)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("got=%#v, want=%#v", err, tc.err)
			}
		})
	}
}

type animal struct {
	Name  string
	Legs  int16
	Tails int8
}

type manimal struct {
	name  string
	legs  int16
	tails int8
}

var (
	animalType = reflect.TypeOf((*animal)(nil)).Elem()
)

func (a manimal) MarshalBoost(w *binser.WBuffer) error {
	w.WriteTypeDescr(animalType) // use same type as animal.
	w.WriteString(a.name)
	w.WriteI16(a.legs)
	w.WriteI8(a.tails)
	return w.Err()
}

func (a *manimal) UnmarshalBoost(r *binser.RBuffer) error {
	r.ReadTypeDescr(animalType) // use same type as animal.
	a.name = r.ReadString()
	a.legs = r.ReadI16()
	a.tails = r.ReadI8()
	return r.Err()
}

var (
	_ binser.Unmarshaler = (*manimal)(nil)
	_ binser.Marshaler   = (*manimal)(nil)
)

func TestRBufferReader(t *testing.T) {
	want := []byte("hello")
	r := binser.NewRBuffer(bytes.NewReader(want))
	got := make([]byte, len(want))
	n, err := r.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	got = got[:n]
	if !bytes.Equal(got, want) {
		t.Fatalf("got=%q, want=%q", got, want)
	}
}

func TestInvalidArray(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := binser.NewEncoder(buf)
	err := enc.Encode([3]int32{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	dec := binser.NewDecoder(buf)
	var v [2]int32
	err = dec.Decode(&v)
	if err == nil {
		t.Fatalf("expected an error!")
	}
	if got, want := err, binser.ErrInvalidArrayLen; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%#v, want=%#v", got, want)
	}
}

func TestDecoderInvalidType(t *testing.T) {
	f, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var iface interface{} = 42

	dec := binser.NewDecoder(f)
	err = dec.Decode(iface)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if got, want := err, binser.ErrTypeNotSupported; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%#v, want=%#v", got, want)
	}
}
