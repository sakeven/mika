package obfs

import (
	"net"
)

type Http struct {
	isServerSide bool
	readStart    bool
	writeStart   bool
	conn         net.Conn
}

func (h *Http) Write(b []byte) (n int, err error) {

	if !h.writeStart {
		if h.isServerSide {
			// write http response header
		} else {
			// write http request header
		}
		h.writeStart = true
	}

	return h.conn.Write(b)
}

func (h *Http) Read(b []byte) (n int, err error) {
	if !h.readStart {
		if h.isServerSide {
			// read http request header
		} else {
			// read http response header
		}
		h.readStart = true
	}

	return h.conn.Read(b)
}

func (h *Http) Close() error {
	return h.conn.Close()
}
