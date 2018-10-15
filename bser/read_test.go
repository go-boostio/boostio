// Copyright 2018 The boostio Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/sbinet/boostio/bser"
)

func TestRead(t *testing.T) {
	type animal struct {
		Name  string
		Legs  int16
		Tails int8
	}

	f, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec, err := bser.NewDecoder(f)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		want interface{}
	}{
		{"bool-false", false},
		{"bool-true", true},
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
		{"string", "hello"},
		{"[3]uint8", [3]uint8{0x11, 0x22, 0x33}},
		{"[]uint8", []uint8{0x11, 0x22, 0x33, 0xff}},
		{"struct", animal{"pet", 4, 1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.New(reflect.TypeOf(tc.want)).Elem()
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
