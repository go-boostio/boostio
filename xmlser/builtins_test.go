// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlser

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func TestBuiltins(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want interface{}
	}{
		{
			raw:  `<v>0</v>`,
			want: false,
		},
		{
			raw:  `<v>1</v>`,
			want: true,
		},
		{
			raw:  `<v><real>2</real><imag>3</imag></v>`,
			want: c64Type{Real: 2, Imag: 3},
		},
		{
			raw:  `<v><real>2</real><imag>3</imag></v>`,
			want: c128Type{Real: 2, Imag: 3},
		},
	} {
		t.Run("", func(t *testing.T) {
			dec := xml.NewDecoder(strings.NewReader(tc.raw))
			val := reflect.New(reflect.TypeOf(tc.want)).Elem()
			err := dec.Decode(val.Addr().Interface())
			if err != nil {
				t.Fatal(err)
			}
			if got, want := val.Interface(), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got=%#v, want=%#v", got, want)
			}
		})
	}
}
