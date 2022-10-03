/**
  Copyright (c) 2022 Zander Schwid & Co. LLC. All rights reserved.
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
	alphabet  [62]byte
	decodeMap [256]byte
	alphabetIdx0 byte
}

// New creates a new base62 encoding.
func New(alphabet []byte) *Encoding {
	enc := &Encoding{}
	copy(enc.alphabet[:], alphabet)
	for i := range enc.decodeMap {
		enc.decodeMap[i] = 255
	}
	for i, b := range enc.alphabet {
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

// Decode decodes a modified base62 string to a byte slice.
func (e * Encoding) DecodeString(b string) []byte {
	answer := big.NewInt(0)
	tmp := new(big.Int)

	for t := b; len(t) > 0; {
		n := len(t)
		if n > 10 {
			n = 10
		}

		total := uint64(0)
		for _, v := range t[:n] {
			ch := e.decodeMap[v]
			if ch == 255 {
				return []byte("")
			}
			total = total*62 + uint64(ch)
		}

		answer.Mul(answer, bigRadix[n])
		tmp.SetUint64(total)
		answer.Add(answer, tmp)

		t = t[n:]
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != e.alphabetIdx0 {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// Encode encodes a byte slice to a modified base62 string.
func  (e * Encoding) EncodeToString(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	maxlen := int(float64(len(b))*1.5) + 1
	answer := make([]byte, 0, maxlen)
	mod := new(big.Int)
	for x.Sign() > 0 {
		x.DivMod(x, bigRadix10, mod)
		if x.Sign() == 0 {
			// When x = 0, we need to ensure we don't add any extra zeros.
			m := mod.Int64()
			for m > 0 {
				answer = append(answer, e.alphabet[m%62])
				m /= 62
			}
		} else {
			m := mod.Int64()
			for i := 0; i < 10; i++ {
				answer = append(answer, e.alphabet[m%62])
				m /= 62
			}
		}
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, e.alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
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
		if i = e.decodeMap[c]; i < 0 {
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
