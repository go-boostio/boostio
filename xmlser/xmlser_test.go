// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlser_test

import (
	"reflect"
)

var typeTestCases = []struct {
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
	{"cmplx64", complex(float32(2), float32(3))},
	{"cmplx128", complex(float64(4), float64(9))},
	{"[3]uint8", [3]uint8{0x11, 0x22, 0x33}},
	{"[]uint8", []uint8{0x11, 0x22, 0x33, 0xff}},
	{"[]byte", []byte("hello")},
	{"string", "hello"},
	{"map[string]string", map[string]string{"eins": "un", "zwei": "deux", "drei": "trois"}},
	{"struct", animal{"pet", 4, 1}},
	{"struct-marshal", manimal{"pet", 4, 1}},
	{"[]string", []string{"s1", "s2", "s3"}},
	{"[]animal", []manimal{{"tiger", 4, 1}, {"monkey", 4, 1}}},
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
