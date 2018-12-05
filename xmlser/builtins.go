// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlser

type c64Type struct {
	Real float32 `xml:"real"`
	Imag float32 `xml:"imag"`
}

type c128Type struct {
	Real float64 `xml:"real"`
	Imag float64 `xml:"imag"`
}
