# base62

[![Base62 CI](https://github.com/arpabet/base62/actions/workflows/build.yaml/badge.svg)](https://github.com/arpabet/base62/actions/workflows/build.yaml)

Canonical, value-preserving Base62 encoding and decoding of byte slices for Go.

The input byte slice is treated as a big-endian arbitrary-precision integer and
converted between radix-256 (bytes) and radix-62 (digits). Each leading zero
byte maps to one leading zero-digit and vice versa, so the encoding is a
**bijection** between byte slices and Base62 strings — `Decode(Encode(b)) == b`
and `Encode(Decode(s)) == s`.

The default alphabet is `0-9a-zA-Z` (digits, then lowercase, then uppercase):
`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`.

## Install

```sh
go get go.arpabet.com/base62
```

Requires Go 1.25+.

## Usage

```go
import "go.arpabet.com/base62"

// Bytes <-> string
s := base62.StdEncoding.EncodeToString(binary)   // []byte -> string
binary, err := base62.StdEncoding.DecodeString(s) // string -> []byte, error

// Unsigned integers (fast path, no big.Int)
s := base62.StdEncoding.EncodeUint64(42)
n, err := base62.StdEncoding.DecodeToUint64(s)
```

### Reusing buffers

For hot paths that want to avoid per-call allocations, use the buffer-based API
together with the length helpers (mirrors `encoding/base64`):

```go
dst := make([]byte, base62.StdEncoding.EncodedLen(len(src)))
n := base62.StdEncoding.Encode(dst, src)        // dst[:n] holds the encoding

out := make([]byte, base62.StdEncoding.DecodedLen(len(dst[:n])))
m, err := base62.StdEncoding.Decode(out, dst[:n]) // out[:m] holds the bytes
```

`EncodedLen(n)` / `DecodedLen(n)` return **upper bounds** on the output size
(the exact length depends on the value, since Base62 is not a fixed-width
encoding). Always slice the result to the returned count.

### Custom alphabet

```go
enc := base62.New([]byte("...62 distinct bytes..."))
```

`New` panics if the alphabet is not exactly 62 bytes or contains duplicates — an
invalid alphabet is a programming error that would otherwise silently corrupt
the decode table.

## Error handling

`DecodeString`, `Decode`, and `DecodeToUint64` return an error (and never panic)
on any byte that is not in the alphabet, including non-ASCII and invalid UTF-8
input. `DecodeToUint64` additionally returns an error on integer overflow.

## CLI

A small command-line tool is provided under `cmd/base62`:

```sh
go build ./cmd/base62

echo "hello world" | base62           # -> 7TqlfhZ 91VHwHy (whitespace preserved per token)
echo "7TqlfhZ 91VHwHy" | base62 -D    # -> hello world
base62 -i input.bin -o output.txt     # file in/out
base62 -v                             # version
```

## Performance notes

Base62 is not a power-of-two radix, so canonical conversion is inherently
**O(n²)** in the input length — every output digit depends on every input byte.
This implementation uses Go's `math/big` (64-bit limb arithmetic) and processes
ten Base62 digits per `uint64` chunk, which is close to optimal for canonical
Base62 in pure Go. It is well suited to the common case of encoding identifiers,
hashes, and keys (tens of bytes); encoding very large blobs (hundreds of KB) is
possible but quadratic.

If you do **not** need value-preserving/canonical output and want O(n) streaming
throughput on large inputs, a bit-grouping Base62 codec (different, longer
output) is a better fit — this library deliberately trades that for canonical
semantics.

### Benchmarks

Representative figures (`go test -bench=.`, Apple M4, Go 1.25). The quadratic
term is visible between 5 KB and 100 KB:

| Operation     | Input  | Throughput | Allocs |
|---------------|--------|------------|--------|
| EncodeToString | 5 KB   | ~5.3 MB/s  | 4      |
| EncodeToString | 100 KB | ~0.26 MB/s | 4      |
| DecodeString   | 5 KB   | ~75 MB/s   | 128    |
| DecodeString   | 100 KB | ~4.7 MB/s  | 2510   |

Small inputs (the common identifier/hash case) run in well under a microsecond:
encoding a 16-byte value takes ~85 ns. Run `make` or `go test -bench=.` to
reproduce on your hardware.

## License

Business Source License 1.1 (BUSL-1.1). See [LICENSE](LICENSE).
