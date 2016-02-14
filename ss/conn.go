package ss

import (
	"io"
	// "log"
	"net"
)

type Conn struct {
	*Crypto
	net.Conn
}

func NewConn(conn net.Conn, crypto *Crypto) *Conn {
	return &Conn{
		Conn:   conn,
		Crypto: crypto,
	}
}

func DailWithRawAddr(network string, rawAddr []byte, server string, cipher *Crypto) (ss net.Conn, err error) {
	conn, err := net.Dial(network, server)
	if err != nil {
		return nil, err
	}

	ss = NewConn(conn, cipher)
	_, err = ss.Write(rawAddr)

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	var buf = make([]byte, 30*1024)
	var cipher = buf

	dataLen := len(b)

	if c.enc == nil {
		var iv []byte
		c.enc, iv = c.newEncStream()
		copy(buf, iv)
		cipher = buf[c.info.ivLen:]
		dataLen += c.info.ivLen
	}

	c.Encrypt(cipher, b)
	return c.Conn.Write(buf[:dataLen])
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.dec == nil {
		iv := make([]byte, c.info.ivLen)
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return 0, err
		}
		c.dec = c.newDecStream(iv)
	}

	buf := make([]byte, 30*1024)
	n, err = c.Conn.Read(buf[:len(b)])
	if err != nil {
		return
	}

	c.Decrypt(b, buf[:n])
	return
}
