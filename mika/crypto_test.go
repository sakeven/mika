package mika

import (
	// "log"
	"reflect"
	"testing"
)

func TestEvpBytesToKey(t *testing.T) {
	key := evpBytesToKey("foobar", 32)
	keyTarget := []byte{0x38, 0x58, 0xf6, 0x22, 0x30, 0xac, 0x3c, 0x91, 0x5f, 0x30, 0x0c, 0x66, 0x43, 0x12, 0xc6, 0x3f, 0x56, 0x83, 0x78, 0x52, 0x96, 0x14, 0xd2, 0x2d, 0xdb, 0x49, 0x23, 0x7d, 0x2f, 0x60, 0xbf, 0xdf}
	if !reflect.DeepEqual(key, keyTarget) {
		t.Errorf("key not correct\n\texpect: %v\n\tgot:   %v\n", keyTarget, key)
	}
}

func Benchmark_evpBytesToKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		evpBytesToKey("foobar", 32)
	}
}

func benchmarkCipherEncInit(b *testing.B, method string) {
	cg := NewCryptoGenerate(method, "password")
	for i := 0; i < b.N; i++ {
		crypto := cg.NewCrypto()
		crypto.initEncStream()
	}
}

func BenchmarkAES256EncInit(b *testing.B) {
	benchmarkCipherEncInit(b, "aes-256-cfb")
}

func benchmarkCipherDecInit(b *testing.B, method string) {
	cg := NewCryptoGenerate(method, "password")
	iv := make([]byte, cg.info.ivLen)
	for i := 0; i < b.N; i++ {
		crypto := cg.NewCrypto()
		crypto.initDecStream(iv)
	}
}

func BenchmarkAES256DecInit(b *testing.B) {
	benchmarkCipherDecInit(b, "aes-256-cfb")
}
