package mika

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_HeaderBytes(t *testing.T) {
	rawAddr := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x10, 0x4}
	h := newHeader(tcpForward, rawAddr)
	h.ChunkId = 123456789

	iv := []byte{0x1, 0x2, 0x3, 0x4}
	key := []byte{0x1, 0x2, 0x3, 0x4}
	bs := h.Bytes(iv, key)
	wanted := []byte{0x1, 0x1, 0x0, 0x0, 0x1, 0x1, 0x2, 0x3, 0x4, 0x5, 0x10, 0x4, 0x0, 0x0, 0x0, 0x0, 0x7, 0x5b, 0xcd, 0x15, 0xc, 0xba, 0x3c, 0xbc, 0x59, 0x66, 0xe, 0xd, 0x9e, 0x84}
	if !reflect.DeepEqual(bs, wanted) {
		t.Errorf("Header convert to bytes error, wanted %#v, got %#v", wanted, bs)
	}
}

func Test_GetHeader(t *testing.T) {
	rawAddr := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x10, 0x4}
	wantedHeader := newHeader(tcpForward, rawAddr)

	iv := []byte{0x1, 0x2, 0x3, 0x4}
	key := []byte{0x1, 0x2, 0x3, 0x4}
	bs := wantedHeader.Bytes(iv, key)

	br := bytes.NewReader(bs)
	gotHeader, err := getHeader(br)
	if err != nil {
		t.Errorf("Parse header error %s", err)
		return
	}

	if gotHeader.ChunkId != wantedHeader.ChunkId {
		t.Errorf("Parse header error, wanted chunk id %#v, got %#v", wantedHeader.ChunkId, gotHeader.ChunkId)
	}
}

func Benchmark_HeaderBytes(b *testing.B) {
	header := newHeader(tcpForward, cipherKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		header.Bytes(cipherIv, cipherKey)
	}
}
