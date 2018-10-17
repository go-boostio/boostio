// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser_test

import (
	"fmt"
	"log"
	"os"

	"github.com/go-boostio/boostio/binser"
)

func ExampleDecoder() {
	f, err := os.Open("testdata/data.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := binser.NewDecoder(f)

	var v1 bool
	err = dec.Decode(&v1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bool: %v\n", v1)

	err = dec.Decode(&v1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bool: %v\n", v1)

	var i8 int8
	err = dec.Decode(&i8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("int8: %#x\n", i8)

	// Output:
	// bool: false
	// bool: true
	// int8: 0x11
}

func ExampleRBuffer() {
	f, err := os.Open("testdata/data.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := binser.NewRBuffer(f)

	fmt.Printf("header: %#v\n", r.ReadHeader())

	fmt.Printf("bool: %v\n", r.ReadBool())
	fmt.Printf("bool: %v\n", r.ReadBool())
	fmt.Printf("int8: %#x\n", r.ReadI8())
	fmt.Printf("int16: %#x\n", r.ReadI16())
	fmt.Printf("int32: %#x\n", r.ReadI32())
	fmt.Printf("int64: %#x\n", r.ReadI64())
	fmt.Printf("uint8: %#x\n", r.ReadU8())
	fmt.Printf("uint16: %#x\n", r.ReadU16())
	fmt.Printf("uint32: %#x\n", r.ReadU32())
	fmt.Printf("uint64: %#x\n", r.ReadU64())
	fmt.Printf("float32: %v\n", r.ReadF32())
	fmt.Printf("float64: %v\n", r.ReadF64())

	// Output:
	// header: binser.Header{Version:0x11, Flags:0x108040804}
	// bool: false
	// bool: true
	// int8: 0x11
	// int16: 0x2222
	// int32: 0x33333333
	// int64: 0x4444444444444444
	// uint8: 0xff
	// uint16: 0x2222
	// uint32: 0x3333333
	// uint64: 0x444444444444444
	// float32: 2.2
	// float64: 3.3
}
