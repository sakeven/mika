package mika

import (
	"fmt"
	"io"
	"net"
)

const (
	socksv5 = 0x05
)

// Socks5TCPRelay as a socks5 server and mika client.
type Socks5TCPRelay struct {
	conn     net.Conn
	cipher   *Crypto
	ssServer string
	closed   bool
}

// NewSocks5TCPRelay creates a new Socks5TCPRelay.
func NewSocks5TCPRelay(conn net.Conn, mikaServer string, cipher *Crypto) *Socks5TCPRelay {
	return &Socks5TCPRelay{
		conn:     conn,
		cipher:   cipher,
		ssServer: mikaServer,
	}
}

// Serve handles connection between socks5 client and remote addr.
func (s *Socks5TCPRelay) Serve() (err error) {
	defer func() {
		if !s.closed {
			s.conn.Close()
		}
	}()
	s.handShake()

	cmd, rawAddr, addr, err := s.parseRequest()
	if err != nil {
		Errorf("Parse request error %v\n", err)
		return
	}

	Infof("Proxy connection to %s\n", string(addr))
	s.reply()

	switch cmd {
	case CONNECT:
		s.connect(rawAddr)
	case UDP_ASSOCIATE:
		s.udpAssociate()
	case BIND:
	default:
		err = fmt.Errorf("unknow cmd type")
	}

	return
}

// version identifier/method selection message
// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+
// reply:
// +----+--------+
// |VER | METHOD |
// +----+--------+
// |  1 |   1    |
// +----+--------+
// handShake dail handshake between socks5 client and socks5 server.
func (s *Socks5TCPRelay) handShake() (err error) {
	raw := make([]byte, 257)
	if _, err = io.ReadFull(s.conn, raw[:2]); err != nil {
		return
	}

	// get socks version
	if ver := raw[0]; ver != socksv5 {
		return fmt.Errorf("error socks version %d", ver)
	}

	// read all method identifier octets
	nmethods := raw[1]
	if _, err = io.ReadFull(s.conn, raw[2:2+nmethods]); err != nil {
		return
	}

	// reply to socks5 client
	_, err = s.conn.Write([]byte{socksv5, 0x00})
	return
}

// The SOCKS request is formed as follows:
//         +----+-----+-------+------+----------+----------+
//         |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
//         +----+-----+-------+------+----------+----------+
//         | 1  |  1  | X’00’ |  1   | Variable |    2     |
//         +----+-----+-------+------+----------+----------+
// Where:
//           o  VER    protocol version: X’05’
//           o  CMD
//              o  CONNECT X’01’
//              o  BIND X’02’
//              o  UDP ASSOCIATE X’03’
//           o  RSV    RESERVED
//           o  ATYP   address type of following address
//              o  IP V4 address: X’01’
//              o  DOMAINNAME: X’03’
//              o  IP V6 address: X’04’
//           o  DST.ADDR       desired destination address
//           o  DST.PORT desired destination port in network octet order

const (
	CONNECT       = 0x01
	BIND          = 0x02
	UDP_ASSOCIATE = 0x03
)

// getCmd gets the cmd requested by socks5 client.
func (s *Socks5TCPRelay) getCmd() (cmd byte, err error) {
	raw := make([]byte, 3)
	if _, err = io.ReadFull(s.conn, raw); err != nil {
		return
	}

	// get socks version
	if ver := raw[0]; ver != socksv5 {
		return 0, fmt.Errorf("error socks version %d", ver)
	}

	cmd = raw[1]
	return
}

// parseRequest parses socks5 client request.
func (s *Socks5TCPRelay) parseRequest() (cmd byte, rawAddr []byte, addr string, err error) {
	cmd, err = s.getCmd()
	if err != nil {
		return
	}

	// check cmd type
	switch cmd {
	case CONNECT:
	case BIND:
	case UDP_ASSOCIATE:
	default:
		err = fmt.Errorf("unknow cmd type")
		return
	}

	if rawAddr, addr, err = getAddress(s.conn); err != nil {
		return
	}

	return
}

// returns a reply formed as follows:
//         +----+-----+-------+------+----------+----------+
//         |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
//         +----+-----+-------+------+----------+----------+
//         | 1  |  1  | X’00’ |  1   | Variable |    2     |
//         +----+-----+-------+------+----------+----------+
// Where:
//           o  VER    protocol version: X’05’
//           o  REP    Reply field:
//              o  X’00’ succeeded
//              o  X’01’ general SOCKS server failure
//              o  X’02’ connection not allowed by ruleset
//              o  X’03’ Network unreachable
//              o  X’04’ Host unreachable
//              o  X’05’ Connection refused
//              o  X’06’ TTL expired
//              o  X’07’ Command not supported
//              o  X’08’ Address type not supported
//              o  X’09’ to X’FF’ unassigned
//           o  RSV    RESERVED
//           o  ATYP   address type of following address
//              o  IP V4 address: X’01’
//              o  DOMAINNAME: X’03’
//              o  IP V6 address: X’04’
//           o  BND.ADDR       server bound address
//           o  BND.PORT       server bound port in network octet order
func (s *Socks5TCPRelay) reply() (err error) {
	_, err = s.conn.Write([]byte{socksv5, 0x00, 0x00, ipv4Addr, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10})
	return
}

// connect handles CONNECT cmd
// Here is a bit magic. It acts as a mika client that redirects conntion to mika server.
func (s *Socks5TCPRelay) connect(rawAddr []byte) (err error) {

	// TODO Dail("tcp", rawAdd) would be more reasonable.
	mika, err := DailWithRawAddr("tcp", s.ssServer, rawAddr, s.cipher)
	if err != nil {
		return
	}

	defer func() {
		if !s.closed {
			err := mika.Close()
			Errorf("Close connection error %v\n", err)
		}
	}()

	go pipe(s.conn, mika)
	pipe(mika, s.conn)
	s.closed = true
	return
}

// udpAssociate handles UDP_ASSOCIATE cmd
func (s *Socks5TCPRelay) udpAssociate() (err error) {
	s.conn.Write([]byte{socksv5, 0x00, 0x00, ipv4Addr, 0x00, 0x00, 0x00, 0x00, 0x04, 0x38})
	return
}
