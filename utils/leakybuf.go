package utils

import (
	"sync"
)

type BufPool struct {
	pool *sync.Pool
	size int
}

func NewBufPool(size int) *BufPool {
	return &BufPool{pool: &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, size)
			return buf
		}},
		size: size,
	}
}

func (bp *BufPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *BufPool) Put(buf []byte) {
	bp.pool.Put(buf)
}

var defaultBufSize = 4096
var leakyBuf = NewBufPool(defaultBufSize)

func GetBuf() []byte {
	return leakyBuf.Get()
}

func PutBuf(buf []byte) {
	leakyBuf.Put(buf)
}
