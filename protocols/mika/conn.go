package mika

import (
	"io"
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/utils"
)

// Conn dails connection between ss server and ss client.
type Conn struct {
	*Crypto
	Conn       protocols.Protocol
	writeBuf   []byte
	readBuf    []byte
	readStart  bool
	writeStart bool
}

// NewConn creates a new shadowsocks connection.
func NewConn(conn protocols.Protocol, crypto *Crypto) protocols.Protocol {
	return &Conn{
		Conn:     conn,
		Crypto:   crypto,
		writeBuf: leakyBuf.Get(),
		readBuf:  leakyBuf.Get(),
	}
}

// RemoteAddr gets remote connection address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Close closes connection and releases buf.
// TODO check close state to avoid close twice.
func (c *Conn) Close() error {
	leakyBuf.Put(c.writeBuf)
	leakyBuf.Put(c.readBuf)
	utils.Debugf("Connection %s closed", c.RemoteAddr())
	return c.Conn.Close()
}

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
		utils.Debugf("Write iv")
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
		utils.Debugf("init enc")
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
		utils.Debugf("Read iv")
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
