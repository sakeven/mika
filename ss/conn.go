package ss

import (
	"io"
	"net"
)

// Conn dails connection between ss server and ss client.
type Conn struct {
	*Crypto
	net.Conn
	writeBuf   []byte
	readBuf    []byte
	readStart  bool
	writeStart bool
}

// NewConn creates a new shadowsocks connection.
func NewConn(conn net.Conn, crypto *Crypto) *Conn {
	return &Conn{
		Conn:     conn,
		Crypto:   crypto,
		writeBuf: leakyBuf.Get(),
		readBuf:  leakyBuf.Get(),
	}
}

// Close closes connection and releases buf.
// TODO check close state to avoid close twice.
func (c *Conn) Close() error {
	leakyBuf.Put(c.writeBuf)
	leakyBuf.Put(c.readBuf)
	Debugf("Connection closed")
	return c.Conn.Close()
}

// func DailWithRawAddr(network string, server string, rawAddr []byte, cipher *Crypto) (ss net.Conn, err error) {
// 	conn, err := net.Dial(network, server)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ss = NewConn(conn, cipher)
// 	_, err = ss.Write(rawAddr)

// 	return
// }

// Write writes data to connection.
func (c *Conn) write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

// Write writes data to connection.
func (c *Conn) Write(b []byte) (n int, err error) {
	// buf = [iv] + [encrypt data]
	// [iv] exists only at beginning of connection, else [iv] is empty.
	var buf = c.writeBuf
	var encryptData = buf

	dataLen := len(b)

	if c.iv == nil {
		Debugf("Write iv")
		c.writeStart = true

		dataLen += c.info.ivLen
		if dataLen > len(buf) {
			buf = make([]byte, dataLen)
		}
		iv := c.initEncStream()
		copy(buf, iv)
		encryptData = buf[c.info.ivLen:]
	} else if !c.writeStart {
		c.initEncStream()
		Debugf("init enc")
		c.writeStart = true
	}

	if dataLen > len(buf) {
		buf = make([]byte, dataLen)
		encryptData = buf
	}

	c.encrypt(encryptData, b)
	return c.Conn.Write(buf[:dataLen])
}

// Read reads data from connection.
func (c *Conn) Read(b []byte) (n int, err error) {
	if c.iv == nil {
		Debugf("Read iv")
		iv := make([]byte, c.info.ivLen)
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return 0, err
		}
		c.initDecStream(iv)
		c.readStart = true
	} else if !c.readStart {
		c.initDecStream(c.iv)
		c.readStart = true
	}

	var buf = c.readBuf
	n, err = c.Conn.Read(buf[:len(b)])
	if err != nil {
		return
	}

	c.decrypt(b, buf[:n])
	return
}
