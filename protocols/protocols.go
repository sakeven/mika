// Package protocols defines protocol interface.
package protocols

import (
	"net"

	"github.com/sakeven/mika/utils"
)

// protocols
const (
	HTTP     = "http"
	SOCKS5   = "socks5"
	KCP      = "kcp"
	TCP      = "tcp"
	ObfsHTTP = "obfs-http"
)

// Protocol in an interface of network connection
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
			// if flusher, ok := dst.(http.Flusher); ok {
			// 	flusher.Flush()
			// }
		}

		if rerr != nil || werr != nil {
			return
		}
	}
}
