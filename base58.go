package base58

import (
	"fmt"
	"math"
	"math/big"
)

const (
	StdPadding = '1'
	NoPadding  = -1

	ripple  = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	bitcoin = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	flickr  = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

// amd64 version in divmod_amd64.s
// 386 version in divmod_386.go
func divmod(a uint64) (q, r uint64)

type CorruptInputError int64

func (c CorruptInputError) Error() string {
	return fmt.Sprintf("illegal base58 data at input index: %d", c)
}

var radix = big.NewInt(58)
var zero = big.NewInt(0)

// Encode encodes x using the encoding enc, writing EncodedLen(len(src)) bytes to dst.
func (e *Encoding) EncodeBigInt(dst []byte, x *big.Int) {
	if x.BitLen() > 64 {
		i := len(dst)
		for x.Cmp(zero) > 0 {
			i--
			mod := new(big.Int)
			x.DivMod(x, radix, mod)
			dst[i] = e.encode[mod.Uint64()]
		}
		e.pad(dst, i)
	} else {
		e.EncodeUint(dst, x.Uint64())
	}
}

// Encode encodes src using the encoding enc, writing EncodedLen(len(src)) bytes to dst.
func (e *Encoding) Encode(dst, src []byte) {
	e.EncodeBigInt(dst, new(big.Int).SetBytes(src))
}

// EncodeUint encodes x using the encoding enc, writing EncodedLen(len(src)) bytes to dst.
func (e *Encoding) EncodeUint(dst []byte, x uint64) {
	i := len(dst)
	var r uint64
	for x > 0 {
		i--
		x, r = divmod(x)
		dst[i] = e.encode[r]
	}
	e.pad(dst, i)
}

func (e *Encoding) pad(dst []byte, i int) {
	if e.padChar != NoPadding {
		for j := 0; dst[j] == 0 && j < i; j++ {
			dst[j] = byte(e.padChar)
		}
	}
}

// EncodeBigIntToString returns the base58 encoding of x.
func (e *Encoding) EncodeBigIntToString(x *big.Int) string {
	dst := make([]byte, EncodedLen(x.BitLen()))
	e.EncodeBigInt(dst, x)
	return string(dst)
}

// EncodeToString returns the base58 encoding of src.
func (e *Encoding) EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)*8))
	e.Encode(dst, src)
	return string(dst)
}

// EncodeUintToString returns the base58 encoding of x.
func (e *Encoding) EncodeUintToString(x uint64) string {
	dst := make([]byte, EncodedUintLen(x))
	e.EncodeUint(dst, x)
	return string(dst)
}

// EncodedLen returns the maximum encoded len in bytes of the base64 encoding
// of src with len n.
func EncodedLen(n int) int {
	return int(math.Ceil(float64(n) / log58))
}

// EncodedLen returns the maximum encoded len in bytes of the base64 encoding
// of n.
func EncodedUintLen(x uint64) int {
	return int(1 + math.Log(float64(x))/log58)
}

func DecodedLen(n int) int {
	return int(math.Floor(float64(n) * log58 / 8))
}

func DecodedUintLen(x uint64) int {
	return int(math.Pow(float64(x), log58)*log58) - 1
}

const log58 = 4.060443010546419

// Encoding is a base58 encoding scheme defined by a 58-character
// alphabet.
type Encoding struct {
	encode    [58]byte
	decodeMap [256]byte
	padChar   rune
}

// WithPadding creates a new encoding identical to enc except
// with a specified padding character, or NoPadding to disable padding.
func (e Encoding) WithPadding(c rune) *Encoding {
	e.padChar = c
	return &e
}

// Ripple is base58 encoding using the Ripple alphabet.
var Ripple = NewEncoding(ripple)

// RawRipple is base58 encoding using the Ripple alphabet without padding.
var RawRipple = Ripple.WithPadding(NoPadding)

// Flickr is base58 encoding using the Flickr alphabet.
var Flickr = NewEncoding(flickr).WithPadding(NoPadding)

// Bitcoin is base58 encoding using the Bitcoin alphabet.
var Bitcoin = NewEncoding(bitcoin)

// RawBitcoin is base58 encoding using the Bitcoin alphabet without padding.
var RawBitcoin = Bitcoin.WithPadding(NoPadding)

// NewEncoding returns a new padded Encoding defined by the given alphabet,
// which must be a 58-byte string. By default padding is turned off.
func NewEncoding(encoder string) *Encoding {
	if len(encoder) != 58 {
		panic("encoding alphabet is not 58-bytes long")
	}

	e := Encoding{padChar: StdPadding}
	copy(e.encode[:], encoder)

	for i := 0; i < len(e.decodeMap); i++ {
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}
	return &e
}

func (e *Encoding) Decode(src []byte) (*big.Int, error) {
	n := new(big.Int)
	j := 0
	for ; j < len(src) && src[j] == byte(e.padChar); j++ {
	}

	for i := range src[j:] {
		c := e.decodeMap[src[i]]
		if c == 0xFF {
			return nil, CorruptInputError(i)
		}
		n.Mul(n, radix)
		n.Add(n, big.NewInt(int64(c)))
	}
	return n, nil
}

func (e *Encoding) DecodeUint(src []byte) (uint64, error) {
	var n uint64
	for i := range src {
		c := e.decodeMap[src[i]]
		if c == 0xFF {
			return 0, CorruptInputError(i)
		}
		n *= 58
		n += uint64(c)
	}
	return n, nil
}
