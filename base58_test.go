package base58

import (
	"math/big"
	"testing"
)

var _s string
var _b *big.Int

func BenchmarkEncodeBigIntToString(b *testing.B) {
	data := make([]byte, 8192)
	data[0] = 0xFF
	b.SetBytes(int64(len(data)))
	n := new(big.Int).SetBytes(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_s = Bitcoin.EncodeBigIntToString(n)
	}
}

func BenchmarkDecodeBigInt(b *testing.B) {
	buf := make([]byte, 8192)
	buf[0] = 0xFF
	data := []byte(Bitcoin.EncodeBigIntToString(new(big.Int).SetBytes(buf)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_b, _ = Bitcoin.Decode(data)
	}
}
