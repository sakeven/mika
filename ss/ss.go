package ss

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	ipv4Addr   = 0x1
	domainAddr = 0x3
	ipv6Addr   = 0x4

	ipv4Len = net.IPv4len
	ipv6Len = net.IPv6len
	portLen = 2
)

// +------+----------+----------+
// | ATYP | DST.ADDR | DST.PORT |
// +------+----------+----------+
// |  1   | Variable |    2     |
// +------+----------+----------+
// o  ATYP    address type of following addresses:
// 		o  IP V4 address: X’01’
// 		o  DOMAINNAME: X’03’
// 		o  IP V6 address: X’04’
// o  DST.ADDR		desired destination address
// o  DST.PORT		desired destination port in network octet
// In an address field (DST.ADDR, BND.ADDR), the ATYP field specifies
//    the type of address contained within the field:
//			o  X’01’
//    the address is a version-4 IP address, with a length of 4 octets
// 			o X’03’
//    the address field contains a fully-qualified domain name.  The first
//    octet of the address field contains the number of octets of name that
//    follow, there is no terminating NUL octet.
//			o  X’04’
//    the address is a version-6 IP address, with a length of 16 octets.
func getAddress(c io.Reader) (raw []byte, addr string, err error) {
	raw = make([]byte, 260)

	pos := 1
	atyp := raw[:pos]
	io.ReadFull(c, atyp)

	var rawAddrLen = 0
	switch atyp[0] {
	case ipv4Addr:
		rawAddrLen = ipv4Len + portLen
	case domainAddr:
		dmLen := raw[pos : pos+1]
		pos++
		io.ReadFull(c, dmLen)
		rawAddrLen = int(dmLen[0] + portLen)
	case ipv6Addr:
		rawAddrLen = ipv6Len + portLen
	default:
		return nil, "", fmt.Errorf("unknow address type %d", atyp[0])
	}

	rawAddr := raw[pos : pos+rawAddrLen]
	io.ReadFull(c, rawAddr)
	pos += rawAddrLen

	var host string
	switch atyp[0] {
	case ipv4Addr:
		host = net.IP(rawAddr[:ipv4Len]).String()
	case domainAddr:
		host = string(rawAddr[:rawAddrLen-portLen])
	case ipv6Addr:
		host = net.IP(rawAddr[:ipv6Len]).String()
	}

	port := int(binary.BigEndian.Uint16(rawAddr[rawAddrLen-portLen:]))

	return raw[:pos], net.JoinHostPort(host, strconv.Itoa(port)), nil
}
