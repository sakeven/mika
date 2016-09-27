package mika

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/utils"

	"github.com/xtaci/kcp-go"
)

// Mika dails connection between mika server and mika client.
type Mika struct {
	*Conn
	header     *header
	serverSide bool
	readBuf    []byte
}

// NewMika wraps a new Mika connection.
// Notice, if header is nil, Mika coonection would be on server side otherwise client side.
func NewMika(conn protocols.Protocol, cipher *Crypto, header *header) (*Mika, error) {
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
		// On server side, we should get header first.
		header, err := getHeader(ss)
		if err != nil {
			return nil, err
		}
		orginHmac := header.Hmac
		header.Bytes(cipher.iv, cipher.key)
		if !bytes.Equal(orginHmac, header.Hmac) {
			return nil, fmt.Errorf("hmac check failed")
		}
		mika.header = header
	} else {
		// On client side, send header as quickly.
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

func DailWithRawAddr(network string, server string, rawAddr []byte, cipher *Crypto) (protocols.Protocol, error) {
	var conn net.Conn
	var err error
	if network == "kcp" {
		var kcpConn *kcp.UDPSession
		kcpConn, err = kcp.DialWithOptions(server, nil, 10, 3)
		kcpConn.SetStreamMode(true)
		kcpConn.SetNoDelay(1, 20, 2, 1)
		kcpConn.SetACKNoDelay(true)
		kcpConn.SetWindowSize(128, 1024)
		conn = kcpConn
	} else {
		conn, err = net.Dial(network, server)
	}

	if err != nil {
		return nil, err
	}

	header := newHeader(tcpForward, rawAddr)
	return NewMika(conn, cipher, header)
}

func DailWithRawAddrHTTP(network string, server string, rawAddr []byte, cipher *Crypto) (protocols.Protocol, error) {
	conn, err := net.Dial(network, server)
	if err != nil {
		return nil, err
	}

	header := newHeader(httpForward, rawAddr)
	return NewMika(conn, cipher, header)
}

// Write writes data to connection.
func (c *Mika) Write(b []byte) (n int, err error) {
	buf := b
	dataLen := len(b)
	if !c.serverSide {
		// ------------------------------
		// | dataLen | hmac | user data |
		// ------------------------------
		// |   2     | 10   | Variable  |
		// ------------------------------
		c.header.ChunkID++
		// hmac := HmacSha1(append(c.iv, c.key...), b)
		// len(hmac)+dataLen = 12
		dataLen += 12

		// utils.Debugf("Send %d hmac %#v", dataLen-12, hmac)
		// utils.Debugf("Data write %#v", b)
		var hmac []byte
		buf, hmac = otaReqChunkAuth(c.iv, c.header.ChunkID, b)
		utils.Debugf("Send %d data, chunkID %d, hmac %#v", dataLen-12, c.header.ChunkID, hmac)
		// utils.Debugf("Data after write %#v", buf[12:dataLen])
	}

	return c.Conn.Write(buf[:dataLen])
}

func (c *Mika) Read(b []byte) (n int, err error) {
	// recover() avoid panic
	var buf = c.readBuf
	// dataLen := len(b)
	if c.serverSide {
		// ------------------------------
		// | dataLen | hmac | user data |
		// ------------------------------
		// |   2     | 10   | Variable  |
		// ------------------------------
		const datePos = 12
		if _, err := io.ReadFull(c.Conn, buf[:datePos]); err != nil {
			return 0, err
		}

		dataLen := int(binary.BigEndian.Uint16(buf[:2]))
		expectedhmac := make([]byte, 10)
		copy(expectedhmac, buf[2:datePos])

		utils.Debugf("dataLen %d expected len %d", dataLen, len(b))

		if dataLen > len(b) {
			utils.Errorf("Date len %d large than b %d", dataLen, len(b))
			return 0, fmt.Errorf("Too large data")
		}

		if _, err := io.ReadFull(c.Conn, b[:dataLen]); err != nil {
			utils.Errorf("Read error %s", err)
			return 0, err
		}

		c.header.ChunkID++
		utils.Debugf("ChunkID %d, receive %d datas, expectedhmac %#v ", c.header.ChunkID, dataLen, expectedhmac)

		_, hmac := otaReqChunkAuth(c.iv, c.header.ChunkID, b[:dataLen])
		if !bytes.Equal(hmac, expectedhmac) {
			utils.Errorf("Hmac %#v mismatch with %#v, remote addr %s", hmac, expectedhmac, c.RemoteAddr())
			return 0, fmt.Errorf("Hmac mismatch")
		}
		// if dataLen > b
		// we should buffer remains.
		return dataLen, nil
	}

	return c.Conn.Read(b)
}
