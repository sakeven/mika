package mika

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"

	"github.com/sakeven/mika/utils"
)

var defaultBufSize = 4096
var leakyBuf = utils.NewBufPool(defaultBufSize)

func HmacSha1(key []byte, data []byte) []byte {
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(data)
	return hmacSha1.Sum(nil)[:10]
}

// TODO use buf to avoid allocate too many memory and objects.
func otaReqChunkAuth(iv []byte, chunkId uint64, data []byte) ([]byte, []byte) {
	nb := make([]byte, 2)
	binary.BigEndian.PutUint16(nb, uint16(len(data)))
	chunkIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(chunkIdBytes, chunkId)
	hmac := HmacSha1(append(iv, chunkIdBytes...), data)
	header := append(nb, hmac...)
	return append(header, data...), hmac
}
