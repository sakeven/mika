package ss

import (
	// "log"
	"net"
)

func pipe(dst, src net.Conn) {
	defer func() {
		dst.Close()
		src.Close()
	}()
	buf := make([]byte, 1024)

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
