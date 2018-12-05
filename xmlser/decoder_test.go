// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-boostio/boostio/xmlser"
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("testdata/data.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := xmlser.NewDecoder(f)
	for _, tc := range typeTestCases[:14] { // FIXME(sbinet): test all test-cases
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
