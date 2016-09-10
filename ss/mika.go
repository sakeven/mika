package ss

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Conn dails connection between ss server and ss client.
type Mika struct {
	*Conn
	header     *header
	serverSide bool
	readBuf    []byte
}

// NewConn wraps a new shadowsocks connection.
// Notice, if header is nil, Conn would be on server side othewise client side.
func NewMika(conn net.Conn, cipher *Crypto, header *header) (*Mika, error) {
	ss := &Conn{
		Conn:     conn,
		Crypto:   cipher,
		writeBuf: leakyBuf.Get(),
		readBuf:  leakyBuf.Get(),
	}

	mika := &Mika{
		Conn:       ss,
		header:     header,
		serverSide: header == nil,
		readBuf:    leakyBuf.Get(),
	}

	if mika.serverSide {
		// on server side, we should get header first.
		header, err := getHeader(ss)
		if err != nil {
			return nil, err
		}
		orginHmac := header.Hamc
		header.Bytes(cipher.iv, cipher.key)
		if !bytes.Equal(orginHmac, header.Hamc) {
			return nil, fmt.Errorf("Hmac check failed")
		}
		mika.header = header
	} else {
		iv := cipher.initEncStream()
		ss.write(iv)
		data := header.Bytes(cipher.iv, cipher.key)
		ss.Write(data)
	}

	return mika, nil
}

// Close closes connection and releases buf.
// TODO check close state to avoid close twice.
func (c *Mika) Close() error {
	leakyBuf.Put(c.readBuf)
	return c.Conn.Close()
}

func DailWithRawAddr(network string, server string, rawAddr []byte, cipher *Crypto) (net.Conn, error) {
	conn, err := net.Dial(network, server)
	if err != nil {
		return nil, err
	}

	header := newHeader(tcpForward, rawAddr)
	return NewMika(conn, cipher, header)
}

// Write writes data to connection.
func (c *Mika) Write(b []byte) (n int, err error) {
	buf := b
	dataLen := len(b)
	if !c.serverSide {
		// ------------------------------
		// | dateLen | hamc | user data |
		// ------------------------------
		// |   2     | 10   | Variable  |
		// ------------------------------
		c.header.ChunkId++
		hmac := HmacSha1(append(c.iv, c.key...), b)
		// len(hmac)+dataLen = 12
		dataLen += 12
		if dataLen > leakyBuf.size {
			buf = make([]byte, dataLen)
		} else {
			buf = leakyBuf.Get()
		}

		Debugf("Send %d data", dataLen-12)
		// Debugf("Send %d hamc %#v", dataLen-12, hmac)
		// Debugf("Data write %#v", b)
		binary.BigEndian.PutUint16(buf, uint16(dataLen-12))
		copy(buf[2:12], hmac)
		copy(buf[12:], b)
		// Debugf("Data after write %#v", buf[12:dataLen])
	}

	return c.Conn.Write(buf[:dataLen])
}

func (c *Mika) Read(b []byte) (n int, err error) {
	// recover() avoid panic
	var buf = c.readBuf
	// dateLen := len(b)
	if c.serverSide {
		// ------------------------------
		// | dateLen | hamc | user data |
		// ------------------------------
		// |   2     | 10   | Variable  |
		// ------------------------------
		const datePos = 12
		if _, err := io.ReadFull(c.Conn, buf[:datePos]); err != nil {
			return 0, err
		}

		dateLen := int(binary.BigEndian.Uint16(buf[:2]))
		expectedhmac := make([]byte, 10)
		copy(expectedhmac, buf[2:datePos])

		Debugf("dataLen %d expected len %d", dateLen, len(b))
		Debugf("expectedhmac %#v ", expectedhmac)

		if dateLen > len(b) {
			Errorf("Date len %d large than b %d", dateLen, len(b))
		}

		if _, err := io.ReadFull(c.Conn, b[:dateLen]); err != nil {
			Errorf("Read error %s", err)
			return 0, err
		}

		hmac := HmacSha1(append(c.iv, c.key...), b[:dateLen])
		if !bytes.Equal(hmac, expectedhmac) {
			Debugf("Hmac %#v mismatch with %#v", hmac, expectedhmac)
			return 0, fmt.Errorf("Hamc mismatch")
		}
		// if dateLen > b
		// we should buffer remains.
		return dateLen, nil
	}

	return c.Conn.Read(b)
}
