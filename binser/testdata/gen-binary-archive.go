// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func main() {
	for _, arch := range []int{
		64,
		32, // needs cross-compilation headers.
	} {
		err := build(arch)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

func build(arch int) error {
	switch arch {
	case 32, 64:
		// ok.
	default:
		return fmt.Errorf("invalid arch %d", arch)
	}

	tmp, err := os.MkdirTemp("", "boostio-binser-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "write.cxx")
	err = os.WriteFile(fname, []byte(src), 0644)
	if err != nil {
		return fmt.Errorf("could not generate C++ source file: %w", err)
	}

	cmd := exec.Command("c++", "-m"+strconv.Itoa(arch), "-lboost_serialization", "-o", "bwrite", "write.cxx")
	cmd.Dir = tmp
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not build C++ Boost program: %w", err)
	}

	archive := new(bytes.Buffer)
	cmd = exec.Command("./bwrite")
	cmd.Dir = tmp
	cmd.Stdout = archive
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run C++ Boost program: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf("testdata/data%d.bin", arch), archive.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("could not save binary archive: %w", err)
	}

	return nil
}

const src = `
#include <boost/archive/binary_oarchive.hpp>
#include <boost/serialization/array.hpp>
#include <boost/serialization/complex.hpp>
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
		: m_name(name)
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

BOOST_CLASS_VERSION(animal, 11)

int main()
{
  binary_oarchive oa{std::cout};

  oa 
	<< false << true
	<< int8_t(0x11)
	<< int16_t(0x2222)
	<< int32_t(0x33333333)
	<< int64_t(0x4444444444444444)
	<< uint8_t(0xff)
	<< uint16_t(0x2222)
	<< uint32_t(0x3333333)
	<< uint64_t(0x444444444444444)
	<< float(2.2)
	<< double(3.3)
	<< std::complex<float>(2.0, 3.0)
	<< std::complex<double>(4.0, 9.0)
	<< std::array<uint8_t, 3>({0x11,0x22,0x33})
	<< std::vector<uint8_t>({0x11,0x22,0x33,0xff})
	<< std::vector<uint8_t>({'h', 'e', 'l', 'l', 'o'})
	<< std::string("hello")
	<< std::map<std::string, std::string>({{"eins", "un"}, {"zwei", "deux"}, {"drei", "trois"}})
	;

  oa << animal("pet", 4, 1);
  oa << animal("pet", 4, 1);
  oa << std::vector<std::string>({"s1", "s2", "s3"});
  oa << std::vector<animal>({animal("tiger",4,1), animal("monkey",4,1)});
}
`
