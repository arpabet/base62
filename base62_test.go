/**
  Copyright (c) 2022 Zander Schwid & Co. LLC. All rights reserved.
*/

package base62_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/schwid/base62"
	"math"
	"math/rand"
	"testing"
)

var stringTests = []struct {
	in  string
	out string
}{
	{"", ""},
	{" ", "w"},
	{"-", "J"},
	{"0", "M"},
	{"1", "N"},
	{"-1", "30B"},
	{"11", "3h7"},
	{"abc", "qMin"},
	{"1234598760", "1a0AFzKIPnihTq"},
	{"abcdefghijklmnopqrstuvwxyz", "hUBXsgd3F2swSlEgbVi2p0Ncr6kzVeJTLaW"},
	{"00000000000000000000000000000000000000000000000000000000000000", "EGCwf6HLNqYIKFfPdd8N0wk949eQseyQb7Rkd652Qk6Akz2Q1ZDjhe3eAAYFYOHESnAVjdMrT9d3FOybe6Y"},
}

var invalidStringTests = []struct {
	in  string
	out string
}{
	{"?", ""},
	{"/", ""},
	{".", ""},
	{"%", ""},
	{"3mJr?", ""},
	{"%3yxU", ""},
	{"3sN#", ""},
	{"4k()", ""},
	{"????", ""},
	{"!@#$%^&*()-_=+~`", ""},
}

var hexTests = []struct {
	in  string
	out string
}{
	{"", ""},
	{"61", "1z"},
	{"626262", "r3lo"},
	{"636363", "rksz"},
	{"73696d706c792061206c6f6e6720737472696e67", "gsYMLccoKcplmYv0sl5XtRVCAdN"},
	{"00eb15231dfceb60925886b67d065299925915aeb172c06647", "02xfEbo02ZLEX6ESUaRlLYJieqVj1OAbB5"},
	{"516b6fcd0f", "69HRUw7"},
	{"bf4f89001e670274dd", "15OLCIkmyVeJD"},
	{"572e4794", "1AZ8hu"},
	{"ecac89cad93923c02321", "5AqnQ3pbuRDGN3"},
	{"10c8511e", "j3pvw"},
	{"00000000000000000000", "0000000000"},
	{"000111d38e5fc9071ffcd20b4a763cc9ae4f252bb4e48fd66a835e252ada93ff480d6dd43dc62a641155a5", "01x2HqU8qh3Dw0z1W2fUcC6mU7O5uQ3DDUZN1Onz3h7rNd4xCsAxeztat"},
	{"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff", "035WzM1EDB9ruSmKv3AmfXhkbYY8j2Am5kR1oXRo0HMobBL8mQlurLTUVcsDuYDTTR0Kdh5ljYUR4AkpUWFSXQz0alF45ZqUwTEh7if5YCQre5MyV5S3IWMe6mYkuhDHaQhaTEhcsGMxKpBHXYDLUujpSzDHMC8jFPX2aKfatfliy11C84eIu86SLYIe7AAEbZqew1Rgh2YJB3rYcofRd2oL1caaMsshz9vFbMjBQEwEV8aWD6qRQf8NdPjq7ikkXlQ81BrpqZXdDY5SEBvihocavXLf0DNPb8Onc2RQ2H7z02p679DgksLv8BwD13MXEgJBvG7l5NXlRzkQrZDN27Z"},
}

func TestBase62(t *testing.T) {
	// Encode tests
	for x, test := range stringTests {
		tmp := []byte(test.in)
		if res := base62.StdEncoding.EncodeToString(tmp); res != test.out {
			t.Errorf("Encode test #%d failed: got: %s want: %s",
				x, res, test.out)
			continue
		} else if rev := base62.StdEncoding.DecodeString(res); !bytes.Equal(tmp, rev) {
			t.Errorf("Decode test #%d failed: got: %q want: %q",
				x, rev, tmp)
			continue
		}
	}

	// Decode tests
	for x, test := range hexTests {
		b, err := hex.DecodeString(test.in)
		if err != nil {
			t.Errorf("hex.DecodeString failed failed #%d: got: %s", x, test.in)
			continue
		}

		if res := base62.StdEncoding.DecodeString(test.out); !bytes.Equal(res, b) {
			t.Errorf("Decode test #%d failed: got: %q want: %q",
				x, res, base62.StdEncoding.EncodeToString(b))
			continue
		}
	}

	// Decode with invalid input
	for x, test := range invalidStringTests {
		if res := base62.StdEncoding.DecodeString(test.in); string(res) != test.out {
			t.Errorf("Decode invalidString test #%d failed: got: %q want: %q",
				x, res, test.out)
			continue
		}
	}
}

func TestEncodeUint64(t *testing.T) {

	s := base62.StdEncoding.EncodeUint64(0)
	if s != "0" {
		t.Errorf("EncodeUint64(%d) = %s, want %s", 0, s, "0")
	}

	for i := 0; i < 100; i++ {
		n := rand.Uint64() % uint64(math.Pow10(i/5))
		actual := base62.StdEncoding.EncodeUint64(n)
		b := marshallUint64(n)
		expected := base62.StdEncoding.EncodeToString(b)
		if actual != expected {
			t.Errorf("EncodeUint64(%d) = %s, want %s", n, actual, expected)
		}
	}

}

func TestDecodeUint64(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := rand.Uint64() % uint64(math.Pow10(i/5))
		src := base62.StdEncoding.EncodeUint64(n)
		got, err := base62.StdEncoding.DecodeToUint64(src)
		if err != nil {
			t.Fatalf("Error occurred while decoding %s (%s).", src, err)
		}
		if got != n {
			t.Errorf("DecodeUint64(%s) = %d, want %d", src, got, n)
		}
	}
}

func TestDecodeUint64Overflow(t *testing.T) {
	src := base62.StdEncoding.EncodeUint64(math.MaxUint64)
	got, err := base62.StdEncoding.DecodeToUint64(src)
	if err != nil {
		t.Fatalf("Error occurred while decoding %s (%s).", src, err)
	}
	if got != math.MaxUint64 {
		t.Errorf("DecodeUint64(%s) = %d, want %d", src, got, uint64(math.MaxUint64))
	}
	bs := []byte(src)
	bs[len(bs)-1]++
	got, err = base62.StdEncoding.DecodeToUint64(string(bs))
	if err == nil {
		t.Errorf("Overflow error should occur while decoding %s but got %d.", bs, got)
	}
	src = "aaaaaaaaaaaaaa"
	got, err = base62.StdEncoding.DecodeToUint64(src)
	if err == nil {
		t.Errorf("Overflow error should occur while decoding %s but got %d.", src, got)
	}
}

func marshallUint64(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return removeLeadingZeros(b)
}

func removeLeadingZeros(b []byte) []byte {
	for i, ch := range b {
		if ch != 0 {
			return b[i:]
		}
	}
	if len(b) > 0 {
		return b[:1]
	}
	return b
}