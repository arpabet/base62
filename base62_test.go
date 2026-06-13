/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package base62_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"go.arpabet.com/base62"
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
		} else if rev, _ := base62.StdEncoding.DecodeString(res); !bytes.Equal(tmp, rev) {
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

		if res, _ := base62.StdEncoding.DecodeString(test.out); !bytes.Equal(res, b) {
			t.Errorf("Decode test #%d failed: got: %q want: %q",
				x, res, base62.StdEncoding.EncodeToString(b))
			continue
		}
	}

	// Decode with invalid input
	for x, test := range invalidStringTests {
		if res, _ := base62.StdEncoding.DecodeString(test.in); string(res) != test.out {
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

// TestDecodeToUint64Invalid verifies that invalid characters are rejected
// rather than silently decoded (the decodeMap sentinel is 255, and the guard
// must compare against it explicitly since the index is an unsigned byte).
func TestDecodeToUint64Invalid(t *testing.T) {
	for _, src := range []string{"!!!", "abc?", "/", "él"} {
		if got, err := base62.StdEncoding.DecodeToUint64(src); err == nil {
			t.Errorf("DecodeToUint64(%q) = %d, want an error", src, got)
		}
	}
}

// TestDecodeStringInvalid verifies that malformed input is rejected with an
// error (and never panics), including non-ASCII bytes and invalid UTF-8.
func TestDecodeStringInvalid(t *testing.T) {
	for _, src := range []string{"abc?", "ab\xffcd", "héllo", "\x80", "????"} {
		got, err := base62.StdEncoding.DecodeString(src)
		if err == nil {
			t.Errorf("DecodeString(%q) = %q, want an error", src, got)
		}
	}
}

// TestDecodeStringNoPanic feeds random/arbitrary bytes through the decoder and
// asserts it returns (cleanly or with an error) without panicking.
func TestDecodeStringNoPanic(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	buf := make([]byte, 64)
	for i := 0; i < 10000; i++ {
		n := r.Intn(len(buf) + 1)
		for j := 0; j < n; j++ {
			buf[j] = byte(r.Intn(256))
		}
		_, _ = base62.StdEncoding.DecodeString(string(buf[:n]))
	}
}

// TestNewInvalidAlphabet verifies that New rejects malformed alphabets instead
// of silently building a corrupt decode table.
func TestNewInvalidAlphabet(t *testing.T) {
	cases := map[string][]byte{
		"empty":      []byte(""),
		"too short":  []byte("0123456789"),
		"too long":   []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0"),
		"duplicates": []byte("00123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY"),
	}
	for name, alphabet := range cases {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("New(%s) did not panic", name)
				}
			}()
			_ = base62.New(alphabet)
		}()
	}
}

// TestBufferAPI checks that the Encode/Decode buffer methods agree with the
// string methods, round-trip cleanly, and stay within the advertised length
// bounds, across edge cases and random inputs (including leading zero bytes).
func TestBufferAPI(t *testing.T) {
	enc := base62.StdEncoding
	r := rand.New(rand.NewSource(2))

	inputs := [][]byte{
		{},
		{0x00},
		{0x00, 0x00, 0x00},
		{0xff},
		{0x00, 0x01, 0x02, 0xff},
	}
	for i := 0; i < 200; i++ {
		b := make([]byte, r.Intn(80))
		for j := range b {
			b[j] = byte(r.Intn(256))
		}
		// Force leading zeros on some inputs.
		if i%3 == 0 && len(b) > 2 {
			b[0], b[1] = 0, 0
		}
		inputs = append(inputs, b)
	}

	for _, in := range inputs {
		// Encode buffer vs string.
		dst := make([]byte, enc.EncodedLen(len(in)))
		n := enc.Encode(dst, in)
		gotStr := enc.EncodeToString(in)
		if string(dst[:n]) != gotStr {
			t.Fatalf("Encode(% x) = %q, EncodeToString = %q", in, dst[:n], gotStr)
		}
		if n > enc.EncodedLen(len(in)) {
			t.Fatalf("Encode wrote %d bytes, exceeds EncodedLen %d", n, enc.EncodedLen(len(in)))
		}

		// Decode buffer vs string, and full round-trip.
		ddst := make([]byte, enc.DecodedLen(n))
		dn, err := enc.Decode(ddst, dst[:n])
		if err != nil {
			t.Fatalf("Decode(%q) error: %v", dst[:n], err)
		}
		if !bytes.Equal(ddst[:dn], in) {
			t.Fatalf("round-trip mismatch: got % x want % x", ddst[:dn], in)
		}
		strDec, err := enc.DecodeString(gotStr)
		if err != nil {
			t.Fatalf("DecodeString(%q) error: %v", gotStr, err)
		}
		if !bytes.Equal(strDec, in) {
			t.Fatalf("DecodeString round-trip mismatch: got % x want % x", strDec, in)
		}
	}
}

// TestDecodeBufferInvalid verifies the buffer Decode reports invalid characters.
func TestDecodeBufferInvalid(t *testing.T) {
	dst := make([]byte, 16)
	if _, err := base62.StdEncoding.Decode(dst, []byte("abc?")); err == nil {
		t.Error("Decode of invalid input returned no error")
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
