package mika

import (
	"fmt"
	"io"
	"net"
)

type Socks5UDPRelay struct {
	conn net.Conn
}

func NewSocks5UDPRelay(conn net.Conn) *Socks5UDPRelay {
	return &Socks5UDPRelay{
		conn: conn,
	}
}

func (s *Socks5UDPRelay) Serve() (err error) {
	defer s.conn.Close()

	rawAddr, _, err := s.parseRequest()

	s.relay(rawAddr)
	return
}

//  +----+------+------+----------+----------+----------+
//  |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
//  +----+------+------+----------+----------+----------+
//  | 2  |  1   |  1   | Variable |    2     | Variable |
//  +----+------+------+----------+----------+----------+
// The fields in the UDP request header are:
//      o  RSV  Reserved X’0000’
//      o  FRAG    Current fragment number
//      o  ATYP    address type of following addresses:
//         o  IP V4 address: X’01’
//         o  DOMAINNAME: X’03’
//         o  IP V6 address: X’04’
// o  DST.ADDR	desired destination address
// o  DST.PORT	desired destination port
// o  DATA	user data
func (s *Socks5UDPRelay) getFrag() (rag byte, err error) {
	raw := make([]byte, 3)
	if _, err = io.ReadFull(s.conn, raw); err != nil {
		return
	}

	return raw[2], nil
}

func (s *Socks5UDPRelay) parseRequest() (rawAddr []byte, addr string, err error) {

	frag, err := s.getFrag()
	if err != nil {
		return
	}

	// An implementation that does not support fragmentation
	// MUST drop any datagram whose FRAG field is other than X’00’.
	if frag != 0x00 {
		return nil, "", fmt.Errorf("frag %d is not 0, drop it", frag)
	}

	return getAddress(s.conn)
}

// relay udp data
func (s *Socks5UDPRelay) relay(rawAddr []byte) (err error) {
	cg := NewCryptoGenerate("aes-128-cfb", "123456")
	cipher := cg.NewCrypto()
	mika, err := DailWithRawAddr("udp", ":8080", rawAddr, cipher)
	if err != nil {
		return
	}
	defer mika.Close()

	go Pipe(s.conn, mika)
	Pipe(mika, s.conn)
	return
}
