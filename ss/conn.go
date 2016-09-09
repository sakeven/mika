package ss

import (
	"io"
	// "log"
	"net"
)

// Conn dails connection between ss server and ss client.
type Conn struct {
	*Crypto
	net.Conn
	writeBuf []byte
	readBuf  []byte
}

func NewConn(conn net.Conn, crypto *Crypto) *Conn {
	return &Conn{
		Conn:     conn,
		Crypto:   crypto,
		writeBuf: leakyBuf.Get(),
		readBuf:  leakyBuf.Get(),
	}
}

func (c *Conn) Close() error {
	leakyBuf.Put(c.writeBuf)
	leakyBuf.Put(c.readBuf)
	Debugf("Connection closed")
	return c.Conn.Close()
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
	// buf = [iv] + [encrypt data]
	// [iv] exists only at beginning of connection, else [iv] is empty.
	// var buf = make([]byte, 30*1024)
	var buf = c.writeBuf
	// 	var buf = leakyBuf.Get()
	// defer leakyBuf.Put(buf)
	var encryptData = buf

	dataLen := len(b)

	if c.enc == nil {
		iv := c.initEncStream()
		copy(buf, iv)
		encryptData = buf[c.info.ivLen:]
		dataLen += c.info.ivLen
	}

	bufLen := len(buf)
	Debugf("dataLen %d bufLen %d", dataLen, bufLen)
	if dataLen > bufLen {
		Errorf("dataLen large than buflen")
	}

	c.encrypt(encryptData, b)
	return c.Conn.Write(buf[:dataLen])
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.dec == nil {
		iv := make([]byte, c.info.ivLen)
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return 0, err
		}
		c.initDecStream(iv)
	}

	// buf := make([]byte, 30*1024)
	var buf = c.readBuf
	// 	var buf = leakyBuf.Get()
	// defer leakyBuf.Put(buf)
	n, err = c.Conn.Read(buf[:len(b)])
	if err != nil {
		return
	}

	c.decrypt(b, buf[:n])
	return
}
