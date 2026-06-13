/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package base62_test

import (
	"bytes"
	"encoding/hex"
	"go.arpabet.com/base62"
	"testing"
)

// FuzzRoundTrip asserts that encoding any byte slice and decoding it back
// returns the original bytes, that the string and buffer APIs agree, and that
// the encoded length never exceeds EncodedLen.
func FuzzRoundTrip(f *testing.F) {
	seeds := [][]byte{
		nil,
		{0x00},
		{0x00, 0x00, 0x00},
		{0xff},
		{0x00, 0x01, 0x02, 0xff},
		[]byte("hello world"),
	}
	for _, h := range hexTests {
		if b, err := hex.DecodeString(h.in); err == nil {
			seeds = append(seeds, b)
		}
	}
	for _, s := range seeds {
		f.Add(s)
	}

	enc := base62.StdEncoding
	f.Fuzz(func(t *testing.T, data []byte) {
		s := enc.EncodeToString(data)
		if len(s) > enc.EncodedLen(len(data)) {
			t.Fatalf("encoded length %d exceeds EncodedLen %d", len(s), enc.EncodedLen(len(data)))
		}

		// String API round-trip.
		got, err := enc.DecodeString(s)
		if err != nil {
			t.Fatalf("DecodeString(%q) error: %v", s, err)
		}
		if !bytes.Equal(got, data) {
			t.Fatalf("round-trip mismatch: got % x want % x", got, data)
		}

		// Buffer API must agree with the string API.
		ebuf := make([]byte, enc.EncodedLen(len(data)))
		en := enc.Encode(ebuf, data)
		if string(ebuf[:en]) != s {
			t.Fatalf("Encode buffer %q != EncodeToString %q", ebuf[:en], s)
		}
		dbuf := make([]byte, enc.DecodedLen(en))
		dn, err := enc.Decode(dbuf, ebuf[:en])
		if err != nil {
			t.Fatalf("Decode buffer error: %v", err)
		}
		if !bytes.Equal(dbuf[:dn], data) {
			t.Fatalf("Decode buffer mismatch: got % x want % x", dbuf[:dn], data)
		}
	})
}

// FuzzDecode asserts that decoding arbitrary strings never panics, and that any
// string which decodes successfully re-encodes to itself (the encoding is a
// bijection between strings and byte slices).
func FuzzDecode(f *testing.F) {
	for _, s := range []string{"", "abc", "?", "ab\xffcd", "0000", "héllo", "\x80"} {
		f.Add(s)
	}
	for _, h := range hexTests {
		f.Add(h.out)
	}

	enc := base62.StdEncoding
	f.Fuzz(func(t *testing.T, s string) {
		out, err := enc.DecodeString(s)
		if err != nil {
			return // invalid input: an error (not a panic) is the contract
		}
		if re := enc.EncodeToString(out); re != s {
			t.Fatalf("non-canonical round-trip: DecodeString(%q) re-encoded to %q", s, re)
		}
		// Buffer decode must agree and stay within DecodedLen.
		dbuf := make([]byte, enc.DecodedLen(len(s)))
		n, err := enc.Decode(dbuf, []byte(s))
		if err != nil {
			t.Fatalf("string decoded but buffer Decode(%q) failed: %v", s, err)
		}
		if !bytes.Equal(dbuf[:n], out) {
			t.Fatalf("buffer Decode disagrees with DecodeString for %q", s)
		}
	})
}
