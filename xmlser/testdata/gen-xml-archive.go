// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	tmp, err := ioutil.TempDir("", "boostio-xmlser-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "write.cxx")
	err = ioutil.WriteFile(fname, []byte(src), 0644)
	if err != nil {
		log.Fatalf("could not generate C++ source file: %v", err)
	}

	cmd := exec.Command("c++", "-lboost_serialization", "-o", "bwrite", "write.cxx")
	cmd.Dir = tmp
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not build C++ Boost program: %v", err)
	}

	archive := new(bytes.Buffer)
	cmd = exec.Command("./bwrite")
	cmd.Dir = tmp
	cmd.Stdout = archive
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not run C++ Boost program: %v", err)
	}

	err = ioutil.WriteFile("testdata/data.xml", archive.Bytes(), 0644)
	if err != nil {
		log.Fatalf("could not save XML archive: %v", err)
	}
}

const src = `
#include <boost/archive/xml_oarchive.hpp>
#include <boost/serialization/array.hpp>
#include <boost/serialization/complex.hpp>
#include <boost/serialization/map.hpp>
#include <boost/serialization/nvp.hpp>
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
		ar & BOOST_SERIALIZATION_NVP(m_name);
		ar & BOOST_SERIALIZATION_NVP(m_legs);
		ar & BOOST_SERIALIZATION_NVP(m_tails);
	}

	std::string m_name;
	int16_t		m_legs;
	int8_t		m_tails;
};

BOOST_CLASS_VERSION(animal, 11)

int main()
{
  xml_oarchive oa{std::cout};

  auto v1 = false;
  auto v2 = true;
  auto v3 = int8_t(0x11);
  auto v4 = int16_t(0x2222);
  auto v5 = int32_t(0x33333333);
  auto v6 = int64_t(0x4444444444444444);
  auto v7 = uint8_t(0xff);
  auto v8 = uint16_t(0x2222);
  auto v9 = uint32_t(0x3333333);
  auto v10 = uint64_t(0x444444444444444);
  auto v11 = float(2.2);
  auto v12 = double(3.3);
  auto v13 = std::complex<float>(2.0, 3.0);
  auto v14 = std::complex<double>(4.0, 9.0);
  auto v15 = std::array<uint8_t, 3>({0x11,0x22,0x33});
  auto v16 = std::vector<uint8_t>({0x11,0x22,0x33,0xff});
  auto v17 = std::vector<uint8_t>({'h', 'e', 'l', 'l', 'o'});
  auto v18 = std::string("hello");
  auto v19 = std::map<std::string, std::string>({{"eins", "un"}, {"zwei", "deux"}, {"drei", "trois"}});

  auto v20 = animal("pet", 4, 1);
  auto v21 = animal("pet", 4, 1);
  auto v22 = std::vector<std::string>({"s1", "s2", "s3"});
  auto v23 = std::vector<animal>({animal("tiger",4,1), animal("monkey",4,1)});

  oa
	<< BOOST_SERIALIZATION_NVP(v1)
	<< BOOST_SERIALIZATION_NVP(v2)
	<< BOOST_SERIALIZATION_NVP(v3)
	<< BOOST_SERIALIZATION_NVP(v4)
	<< BOOST_SERIALIZATION_NVP(v5)
	<< BOOST_SERIALIZATION_NVP(v6)
	<< BOOST_SERIALIZATION_NVP(v7)
	<< BOOST_SERIALIZATION_NVP(v8)
	<< BOOST_SERIALIZATION_NVP(v9)
	<< BOOST_SERIALIZATION_NVP(v10)
	<< BOOST_SERIALIZATION_NVP(v11)
	<< BOOST_SERIALIZATION_NVP(v12)
	<< BOOST_SERIALIZATION_NVP(v13)
	<< BOOST_SERIALIZATION_NVP(v14)
	<< BOOST_SERIALIZATION_NVP(v15)
	<< BOOST_SERIALIZATION_NVP(v16)
	<< BOOST_SERIALIZATION_NVP(v17)
	<< BOOST_SERIALIZATION_NVP(v18)
	<< BOOST_SERIALIZATION_NVP(v19)
	<< BOOST_SERIALIZATION_NVP(v20)
	<< BOOST_SERIALIZATION_NVP(v21)
	<< BOOST_SERIALIZATION_NVP(v22)
	<< BOOST_SERIALIZATION_NVP(v23)
	;
}
`
