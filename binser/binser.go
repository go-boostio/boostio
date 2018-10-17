// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binser provides types to read and write binary archives from the C++
// Boost Serialization library.
package binser // import "github.com/go-boostio/boostio/binser"

//go:generate go run ./testdata/gen-binary-archive.go

import (
	"github.com/pkg/errors"
)

var (
	errNotBoost = errors.New("binser: not a Boost binary archive")
)

// Header describes a binary boost archive.
type Header struct {
	Version uint16
	Flags   uint64
}

// Unmarshaler is the interface implemented by types that can unmarshal a binary
// Boost description of themselves.
type Unmarshaler interface {
	UnmarshalBoost(r *Reader) error
}
