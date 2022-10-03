/**
  Copyright (c) 2022 Zander Schwid & Co. LLC. All rights reserved.
*/

package base62_test

import (
	"bytes"
	"github.com/schwid/base62"
	"testing"
)

var (
	raw5k       = bytes.Repeat([]byte{0xff}, 5000)
	raw100k     = bytes.Repeat([]byte{0xff}, 100*1000)
	encoded5k   = base62.StdEncoding.EncodeToString(raw5k)
	encoded100k = base62.StdEncoding.EncodeToString(raw100k)
)

func BenchmarkBase62Encode_5K(b *testing.B) {
	b.SetBytes(int64(len(raw5k)))
	for i := 0; i < b.N; i++ {
		base62.StdEncoding.EncodeToString(raw5k)
	}
}

func BenchmarkBase62Encode_100K(b *testing.B) {
	b.SetBytes(int64(len(raw100k)))
	for i := 0; i < b.N; i++ {
		base62.StdEncoding.EncodeToString(raw100k)
	}
}

func BenchmarkBase62Decode_5K(b *testing.B) {
	b.SetBytes(int64(len(encoded5k)))
	for i := 0; i < b.N; i++ {
		base62.StdEncoding.DecodeString(encoded5k)
	}
}

func BenchmarkBase62Decode_100K(b *testing.B) {
	b.SetBytes(int64(len(encoded100k)))
	for i := 0; i < b.N; i++ {
		base62.StdEncoding.DecodeString(encoded100k)
	}
}
