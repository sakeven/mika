package mika

import (
	"net"
)

func Pipe(dst, src net.Conn) {
	var buf = leakyBuf.Get()

	defer func() {
		leakyBuf.Put(buf)
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
