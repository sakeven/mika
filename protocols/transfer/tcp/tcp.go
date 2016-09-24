package tcp

import (
	"net"
	"time"
)

type Conn struct {
	net.Conn
	Timeout time.Duration
}

func (c *Conn) Write(b []byte) (int, error) {
	c.SetWriteDeadline(time.Now().Add(c.Timeout))
	return c.Conn.Write(b)
}

func (c *Conn) Read(b []byte) (int, error) {
	c.SetReadDeadline(time.Now().Add(c.Timeout))
	return c.Conn.Read(b)
}
