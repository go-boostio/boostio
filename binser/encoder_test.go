// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"github.com/go-boostio/boostio/binser"
)

var encoderTestCases = []struct {
	name string
	want interface{}
}{
	{"bool-true", true},
	{"bool-false", false},
	{"int8", int8(0x11)},
	{"int16", int16(0x2222)},
	{"int32", int32(0x33333333)},
	{"int64", int64(0x4444444444444444)},
	{"uint8", uint8(0xff)},
	{"uint16", uint16(0x2222)},
	{"uint32", uint32(0x3333333)},
	{"uint64", uint64(0x444444444444444)},
	{"float32", float32(2.2)},
	{"float64", 3.3},
	{"[3]uint8", [3]uint8{0x11, 0x22, 0x33}},
	{"[]uint8", []uint8{0x11, 0x22, 0x33, 0xff}},
	{"[]byte", []byte("hello")},
	{"string", "hello"},
	{"map[string]string", map[string]string{"eins": "un", "zwei": "deux", "drei": "trois"}},
	{"struct", animal{"pet", 4, 1}},
	{"struct-marshal", manimal{"pet", 4, 1}},
}

func TestEncoder(t *testing.T) {
	type animal struct {
		Name  string
		Legs  int16
		Tails int8
	}

	for _, tc := range encoderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buf = new(bytes.Buffer)
				err error
				got = reflect.New(reflect.TypeOf(tc.want)).Elem()
			)

			enc := binser.NewEncoder(buf)
			err = enc.Encode(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got.Kind() == reflect.Map {
				got.Set(reflect.MakeMap(got.Type()))
			}

			dec := binser.NewDecoder(bytes.NewReader(buf.Bytes()))
			err = dec.Decode(got.Addr().Interface())
			if err != nil {
				t.Fatalf("could not decode value: %v\n%s", err, hex.Dump(buf.Bytes()))
			}

			if got, want := got.Interface(), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed:\ngot= %#v (%T)\nwant=%#v (%T)", got, got, want, want)
			}
		})
	}
}

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestEncoderError(t *testing.T) {
	for _, tc := range encoderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := binser.NewEncoder(errWriter{})
			err := enc.Encode(tc.want)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	}
}

func TestWBufferWriter(t *testing.T) {
	want := []byte("hello")
	buf := new(bytes.Buffer)
	w := binser.NewWBuffer(buf)
	_, err := w.Write(want)
	if err != nil {
		t.Fatal(err)
	}

	if got := buf.Bytes(); !bytes.Equal(got, want) {
		t.Fatalf("got=%q, want=%q", got, want)
	}
}
