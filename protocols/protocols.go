// Package protocols defines protocol interface.
package protocols

import (
	"net"

	"github.com/sakeven/mika/utils"
)

type Protocol interface {
	Write(b []byte) (n int, err error)
	Read(b []byte) (n int, err error)
	RemoteAddr() net.Addr
	Close() error
}

// Pipe pipes two protocols which should be in same layer.
func Pipe(dst, src Protocol) {
	var buf = utils.GetBuf()

	defer func() {
		utils.PutBuf(buf)
		dst.Close()
	}()

	var rerr, werr error
	var n int
	for {
		n, rerr = src.Read(buf)
		if n > 0 {
			_, werr = dst.Write(buf[:n])
		}

		if rerr != nil || werr != nil {
			return
		}
	}
}
