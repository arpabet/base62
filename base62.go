/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package base62

import (
	"fmt"
	"math/big"
)

const (
	radix = uint64(62)
)

type Encoding struct {
	alphabet     [62]byte
	decodeMap    [256]byte
	alphabetIdx0 byte
}

// New creates a new base62 encoding from the given alphabet.
//
// The alphabet must be exactly 62 bytes long and contain no duplicate bytes;
// New panics otherwise, since an invalid alphabet is a programming error that
// would silently corrupt the decode table.
func New(alphabet []byte) *Encoding {
	if len(alphabet) != 62 {
		panic(fmt.Sprintf("base62: alphabet must be 62 bytes, got %d", len(alphabet)))
	}
	enc := &Encoding{}
	copy(enc.alphabet[:], alphabet)
	for i := range enc.decodeMap {
		enc.decodeMap[i] = 255
	}
	for i, b := range enc.alphabet {
		if enc.decodeMap[b] != 255 {
			panic(fmt.Sprintf("base62: duplicate byte %#x at index %d in alphabet", b, i))
		}
		enc.decodeMap[b] = byte(i)
	}
	enc.alphabetIdx0 = alphabet[0]
	return enc
}

var StdEncoding = New([]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))

var bigRadix = [...]*big.Int{
	big.NewInt(0),
	big.NewInt(62),
	big.NewInt(62 * 62),
	big.NewInt(62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62 * 62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62 * 62 * 62 * 62 * 62),
	big.NewInt(62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62),
	bigRadix10,
}

var bigRadix10 = big.NewInt(62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62) // 62^10

// DecodeString decodes a modified base62 string to a byte slice.
func (e *Encoding) DecodeString(b string) ([]byte, error) {
	answer := big.NewInt(0)
	tmp := new(big.Int)

	for t := b; len(t) > 0; {
		n := len(t)
		if n > 10 {
			n = 10
		}

		total := uint64(0)
		for k := 0; k < n; k++ {
			v := t[k]
			c := e.decodeMap[v]
			if c == 255 {
				return nil, fmt.Errorf("invalid character '%c' in decoding a base62 string '%s'", v, b)
			}
			total = total*62 + uint64(c)
		}

		answer.Mul(answer, bigRadix[n])
		tmp.SetUint64(total)
		answer.Add(answer, tmp)

		t = t[n:]
	}

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != e.alphabetIdx0 {
			break
		}
	}
	magLen := (answer.BitLen() + 7) / 8
	val := make([]byte, numZeros+magLen)
	if magLen > 0 {
		answer.FillBytes(val[numZeros:])
	}

	return val, nil
}

// EncodeToString encodes a byte slice to a modified base62 string.
func (e *Encoding) EncodeToString(b []byte) string {
	dst := make([]byte, e.EncodedLen(len(b)))
	n := e.Encode(dst, b)
	return string(dst[:n])
}

// EncodedLen returns the maximum number of bytes the base62 encoding of an
// n-byte input may occupy. The exact length depends on the input value (it is
// shorter when the leading bytes are zero); use this to size the destination
// buffer for Encode.
func (e *Encoding) EncodedLen(n int) int {
	// ceil(n * log(256)/log(62)); log(256)/log(62) ~= 1.3436 < 1.346.
	return n*1346/1000 + 2
}

// DecodedLen returns the maximum number of bytes the decoding of an n-character
// base62 string may occupy. The worst case (and exact length for an input of
// all zero-digits) is n; the result is shorter for inputs whose leading digits
// are non-zero. Use this to size the destination buffer for Decode.
func (e *Encoding) DecodedLen(n int) int {
	return n
}

// Encode encodes src into base62 and writes the result to dst, returning the
// number of bytes written. dst must be at least EncodedLen(len(src)) bytes
// long. The encoded bytes are written left-aligned at dst[0].
func (e *Encoding) Encode(dst, src []byte) int {
	x := new(big.Int).SetBytes(src)
	mod := new(big.Int)

	// Emit digits least-significant first into the tail of dst, then
	// left-align, avoiding a separate reversal pass.
	pos := len(dst)
	for x.Sign() > 0 {
		x.DivMod(x, bigRadix10, mod)
		m := mod.Int64()
		if x.Sign() == 0 {
			// Most-significant chunk: no leading-zero padding.
			for m > 0 {
				pos--
				dst[pos] = e.alphabet[m%62]
				m /= 62
			}
		} else {
			for i := 0; i < 10; i++ {
				pos--
				dst[pos] = e.alphabet[m%62]
				m /= 62
			}
		}
	}

	// Leading zero bytes become leading zero-digits.
	for _, b := range src {
		if b != 0 {
			break
		}
		pos--
		dst[pos] = e.alphabetIdx0
	}

	n := len(dst) - pos
	copy(dst, dst[pos:])
	return n
}

// Decode decodes the base62 bytes in src and writes the result to dst,
// returning the number of bytes written. dst must be at least
// DecodedLen(len(src)) bytes long. It returns an error on an invalid character.
func (e *Encoding) Decode(dst, src []byte) (int, error) {
	answer := big.NewInt(0)
	tmp := new(big.Int)

	for t := src; len(t) > 0; {
		n := len(t)
		if n > 10 {
			n = 10
		}

		total := uint64(0)
		for k := 0; k < n; k++ {
			c := e.decodeMap[t[k]]
			if c == 255 {
				return 0, fmt.Errorf("invalid character '%c' in decoding a base62 string '%s'", t[k], src)
			}
			total = total*62 + uint64(c)
		}

		answer.Mul(answer, bigRadix[n])
		tmp.SetUint64(total)
		answer.Add(answer, tmp)

		t = t[n:]
	}

	var numZeros int
	for numZeros = 0; numZeros < len(src); numZeros++ {
		if src[numZeros] != e.alphabetIdx0 {
			break
		}
	}
	magLen := (answer.BitLen() + 7) / 8
	n := numZeros + magLen
	for i := 0; i < numZeros; i++ {
		dst[i] = 0
	}
	if magLen > 0 {
		answer.FillBytes(dst[numZeros:n])
	}
	return n, nil
}

// EncodeUint64 encodes the unsigned integer.
func (e *Encoding) EncodeUint64(n uint64) string {
	if n == 0 {
		return string(e.alphabetIdx0)
	}
	answer := make([]byte, 12)
	i := len(answer)
	var mod uint64
	for n > 0 {
		n, mod = n/radix, n%radix
		i--
		answer[i] = e.alphabet[mod]
	}
	return string(answer[i:])
}

// DecodeUint64 decodes the base62 encoded string to an unsigned integer.
func (e *Encoding) DecodeToUint64(src string) (uint64, error) {
	var n, m uint64
	var i byte
	for _, c := range []byte(src) {
		if i = e.decodeMap[c]; i == 255 {
			return 0, fmt.Errorf("invalid character '%c' in decoding a base62 string %q", c, src)
		}
		m = n*radix + uint64(i)
		if m < n {
			return 0, fmt.Errorf("overflow in decoding a base62 string %q", src)
		}
		n = m
	}
	return n, nil
}
