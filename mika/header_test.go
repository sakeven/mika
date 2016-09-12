package mika

import (
	"testing"
)

func Test_HeaderBytes(t *testing.T) {

}

func Test_GetHeader(t *testing.T) {

}

func Benchmark_HeaderBytes(b *testing.B) {
	header := newHeader(tcpForward, cipherKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		header.Bytes(cipherIv, cipherKey)
	}
}
