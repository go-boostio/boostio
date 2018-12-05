// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package boostio provides the general infrastructure to read and write
// streams compatible with the C++ Boost Serialization library:
//  - https://theboostcpplibraries.com/boost.serialization
package boostio // import "github.com/go-boostio/boostio"

const (
	Version uint16 = 0x11 // Boost archive version
)
