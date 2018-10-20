// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-boostio/boostio/binser"
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

func TestEncoder(t *testing.T) {
	type animal struct {
		Name  string
		Legs  int16
		Tails int8
	}

	for _, tc := range typeTestCases {
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
	for _, tc := range typeTestCases {
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

func TestEncoderInvalidType(t *testing.T) {
	var iface interface{} = 42

	enc := binser.NewEncoder(new(bytes.Buffer))
	err := enc.Encode(iface)
	if err == nil {
		t.Fatalf("expected an error")
	}

	if got, want := err, binser.ErrTypeNotSupported; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%#v, want=%#v", got, want)
	}
}

func TestEncoderCompatWithBoost(t *testing.T) {
	f, err := os.Create("testdata/check.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	enc := binser.NewEncoder(f)
	for _, tc := range typeTestCases {
		err := enc.Encode(tc.want)
		if err != nil {
			t.Fatalf("error encoding %q: %v", tc.name, err)
		}
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("error closing output stream: %v", err)
	}

	tmp, err := ioutil.TempDir("", "boostio-binser-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "read.cxx")
	err = ioutil.WriteFile(fname, []byte(boostReadSrc), 0644)
	if err != nil {
		log.Fatalf("could not generate C++ source file: %v", err)
	}

	dbg := new(bytes.Buffer)
	cmd := exec.Command("c++", "-std=c++11", "-lboost_serialization", "-o", "bread", "read.cxx")
	cmd.Dir = tmp
	cmd.Stdout = dbg
	cmd.Stderr = dbg
	err = cmd.Run()
	if err != nil {
		t.Skipf("could not compile C++ Boost: %s", dbg.Bytes())
		os.Remove(f.Name())
		return
	}

	archive, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	out := new(bytes.Buffer)
	cmd = exec.Command(filepath.Join(tmp, "bread"))
	cmd.Stdin = bytes.NewReader(archive)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error reading back boost archive: %v\n%s", err, out.Bytes())
	}
	want := `bool: 0
bool: 1
int8_t: 0x11
int16_t: 0x2222
int32_t: 0x33333333
int64_t: 0x44444444
uint8_t: 0xff
uint16_t: 0x2222
uint32_t: 0x3333333
uint64_t: 0x44444444
float32: 2.2
float64: 3.3
[3]uint8: {0x11, 0x22, 0x33, }
[]uint8: {0x11, 0x22, 0x33, 0xff, }
[]uint8: {68, 65, 6c, 6c, 6f, }
string: "hello"
map: {{drei: trois}, {eins: un}, {zwei: deux}, }
animal: {name: pet, legs: 4, tails: 1}
animal: {name: pet, legs: 4, tails: 1}
[]string: {s1, s2, s3, }
[]animal: {{name: tiger, legs: 4, tails: 1}, {name: monkey, legs: 4, tails: 1}, }
`
	if got, want := out.Bytes(), []byte(want); !bytes.Equal(got, want) {
		t.Fatalf("output differs:\ngot:\n%s\nwant:%s\n", got, want)
	}

	os.Remove(f.Name())
}

const boostReadSrc = `
#include <boost/archive/binary_iarchive.hpp>
#include <boost/serialization/array.hpp>
#include <boost/serialization/map.hpp>
#include <boost/serialization/vector.hpp>

#include <iostream>
#include <string>
#include <vector>
#include <array>

#include <stdint.h>

using namespace boost::archive;

class animal {
public:
	animal(std::string name = "pet", int legs=4, int tails=2) 
		: m_name("pet")
		, m_legs(legs)
		, m_tails(tails)
	{}

	std::string name()  const { return m_name; }
	int			legs()  const { return m_legs; }
	int			tails() const { return m_tails; }

private:

	friend class boost::serialization::access;

	template <typename Archive>
	void serialize(Archive &ar, const unsigned int version) {
		ar & m_name;
		ar & m_legs;
		ar & m_tails;
	}

	std::string m_name;
	int16_t		m_legs;
	int8_t		m_tails;
};

int main()
{
  binary_iarchive ia{std::cin};

  {
	bool v;
	ia >> v;
	std::cout << "bool: " << v << "\n";
  }

  {
	bool v;
	ia >> v;
	std::cout << "bool: " << v << "\n";
  }

  {
	int8_t v;
	ia >> v;
	std::printf("int8_t: 0x%x\n", v);
  }

  {
	int16_t v;
	ia >> v;
	std::printf("int16_t: 0x%x\n", v);
  }

  {
	int32_t v;
	ia >> v;
	std::printf("int32_t: 0x%x\n", v);
  }

  {
	int64_t v;
	ia >> v;
	std::printf("int64_t: 0x%x\n", v);
  }

  {
	uint8_t v;
	ia >> v;
	std::printf("uint8_t: 0x%x\n", v);
  }

  {
	uint16_t v;
	ia >> v;
	std::printf("uint16_t: 0x%x\n", v);
  }

  {
	uint32_t v;
	ia >> v;
	std::printf("uint32_t: 0x%x\n", v);
  }

  {
	uint64_t v;
	ia >> v;
	std::printf("uint64_t: 0x%x\n", v);
  }

  {
	float v;
	ia >> v;
	std::printf("float32: %1.1f\n", v);
  }

  {
	double v;
	ia >> v;
	std::printf("float64: %1.1f\n", v);
  }

  {
    std::array<uint8_t, 3> v;
	ia >> v;
	std::cout << "[3]uint8: {";
	for (auto i : v) { std::printf("0x%x, ", i); }
	std::cout << "}\n";
  }

  {
    std::vector<uint8_t> v;
	ia >> v;
	std::cout << "[]uint8: {";
	for (auto i : v) { std::printf("0x%x, ", i); }
	std::cout << "}\n";
  }

  {
    std::vector<uint8_t> v;
	ia >> v;
	std::cout << "[]uint8: {";
	for (auto i : v) { std::printf("%x, ", i); }
	std::cout << "}\n";
  }

  {
    std::string v;
	ia >> v;
	std::cout << "string: \"" << v << "\"\n";
  }

  {
	std::map<std::string, std::string> v;
	ia >> v;
	std::cout << "map: {";
	for (const auto &kv : v) { std::cout << "{" <<kv.first << ": " << kv.second << "}, "; } 
	std::cout << "}\n";
  }

  {
	  animal v;
	  ia >> v;
	  std::cout << "animal: {name: " << v.name() << ", legs: " << v.legs() << ", tails: " << v.tails() << "}\n";
  }

  {
	  animal v;
	  ia >> v;
	  std::cout << "animal: {name: " << v.name() << ", legs: " << v.legs() << ", tails: " << v.tails() << "}\n";
  }

  {
	  std::vector<std::string> vs;
	  ia >> vs;
	  std::cout << "[]string: {";
	  for (auto v: vs) { std::cout << v << ", "; }
	  std::cout << "}\n";
  }

  {
	  std::vector<animal> vs;
	  ia >> vs;
	  std::cout << "[]animal: {";
	  for (auto v: vs) { std::cout << "{name: " << v.name() << ", legs: " << v.legs() << ", tails: " << v.tails() << "}, "; }
	  std::cout << "}\n";
  }
}
`
