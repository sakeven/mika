package ss

import (
	// "log"
	"net"
)

func pipe(dst, src net.Conn) {
	var buf = leakyBuf.Get()

	defer func() {
		leakyBuf.Put(buf)
		dst.Close()
	}()

	// buf := make([]byte, 4096)

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
