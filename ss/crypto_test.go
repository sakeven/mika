package ss

import (
	"crypto/cipher"
	"crypto/rand"
	"io"
	// "log"
	"testing"
)

func Test_evpBytesToKey(t *testing.T) {
	b := evpBytesToKey("123456", 16)
	dest := []byte{225, 10, 220, 57, 73, 186, 89, 171, 190, 86, 224, 87, 242, 15, 136, 62}
	if string(b) != string(dest) {
		t.Errorf("Get error")
	}
}

var cipherKey = make([]byte, 64)
var cipherIv = make([]byte, 64)

func init() {
	for i := 0; i < len(cipherKey); i++ {
		cipherKey[i] = byte(i)
	}
	io.ReadFull(rand.Reader, cipherIv)
}

func benchmarkCipherInit(b *testing.B, method string) {
	ci := cryptoInfoMap[method]
	key := cipherKey[:ci.keyLen]
	block := ci.newBlock(key)

	iv := make([]byte, ci.ivLen)
	for i := 0; i < b.N; i++ {
		cipher.NewCFBEncrypter(block, iv)
	}
}

func BenchmarkAES256Init(b *testing.B) {
	benchmarkCipherInit(b, "aes-256-cfb")
}
