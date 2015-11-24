package base58

import (
	"math/big"
	"testing"
)

var _s string
var _b *big.Int

func expect(t *testing.T, want, got string) {
	if want != got {
		t.Fatalf("want: %q, got %q", want, got)
	}
}

func TestEncode(t *testing.T) {
	testStr := "this is a test"
	expect(t, "jo91waLQA1NNeBmZKUF", Bitcoin.EncodeToString([]byte(testStr)))
}

func TestDecode(t *testing.T) {
	testStr := "this is a test"
	enc := Bitcoin.EncodeToString([]byte(testStr))
	expect(t, "jo91waLQA1NNeBmZKUF", enc)
	i, err := Bitcoin.Decode([]byte(enc))
	if err != nil {
		t.Fatal(err)
	}
	dec := string(i.Bytes())
	expect(t, testStr, dec)
}

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

func BenchmarkEncodeBigIntToStringSmall(b *testing.B) {
	data := make([]byte, 8)
	data[0] = 0xFF
	b.SetBytes(int64(len(data)))
	n := new(big.Int).SetBytes(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_s = Bitcoin.EncodeBigIntToString(n)
	}
}

func BenchmarkDecodeBigIntSmall(b *testing.B) {
	buf := make([]byte, 8)
	buf[0] = 0xFF
	data := []byte(Bitcoin.EncodeBigIntToString(new(big.Int).SetBytes(buf)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_b, _ = Bitcoin.Decode(data)
	}
}
