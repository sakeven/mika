package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Addr type and length
const (
	AddrIPv4   = 0x1
	AddrDomain = 0x3
	AddrIPv6   = 0x4

	ipv4Len = net.IPv4len
	ipv6Len = net.IPv6len
	portLen = 2
)

// GetAddress gets address from a reader.
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
func GetAddress(c io.Reader) (raw []byte, addr string, err error) {
	raw = make([]byte, 260)

	pos := 1
	atyp := raw[:pos]
	io.ReadFull(c, atyp)
	var errs []error

	var rawAddrLen = 0
	switch atyp[0] {
	case AddrIPv4:
		rawAddrLen = ipv4Len + portLen
	case AddrDomain:
		dmLen := raw[pos : pos+1]
		pos++
		io.ReadFull(c, dmLen)
		rawAddrLen = int(dmLen[0] + portLen)
	case AddrIPv6:
		rawAddrLen = ipv6Len + portLen
	default:
		errs = append(errs, fmt.Errorf("unknow address type %d", atyp[0]))
		//treat ad domain
		dmLen := raw[pos : pos+1]
		pos++
		io.ReadFull(c, dmLen)
		rawAddrLen = int(dmLen[0] + portLen)
	}

	rawAddr := raw[pos : pos+rawAddrLen]
	io.ReadFull(c, rawAddr)
	pos += rawAddrLen

	var host string
	switch atyp[0] {
	case AddrIPv4:
		host = net.IP(rawAddr[:ipv4Len]).String()
	case AddrDomain:
		host = string(rawAddr[:rawAddrLen-portLen])
	case AddrIPv6:
		host = net.IP(rawAddr[:ipv6Len]).String()
	}

	port := int(binary.BigEndian.Uint16(rawAddr[rawAddrLen-portLen:]))

	if len(errs) > 0 {
		return nil, "", errs[0]
	}

	return raw[:pos], net.JoinHostPort(host, strconv.Itoa(port)), nil
}

// ToAddr changes addr from string to network bytes
func ToAddr(addr string) []byte {
	if strings.Index(addr, ":") < 0 {
		addr += ":80"
	}

	host, port, err := net.SplitHostPort(addr) //stats.g.doubleclick.net:443
	if err != nil {
		return nil
	}
	addrBytes := make([]byte, 0)
	ip := net.ParseIP(host)

	if ip == nil {
		l := len(host)
		addrBytes = append(addrBytes, AddrDomain)
		addrBytes = append(addrBytes, byte(l))
		addrBytes = append(addrBytes, []byte(host)...)
	} else if ipv4 := ip.To4(); ipv4 != nil {
		addrBytes = append(addrBytes, AddrIPv4)
		addrBytes = append(addrBytes, []byte(ipv4)...)
	} else {
		addrBytes = append(addrBytes, AddrIPv6)
		addrBytes = append(addrBytes, []byte(ip)...)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil
	}

	bp := make([]byte, 2)
	binary.BigEndian.PutUint16(bp, uint16(p))

	addrBytes = append(addrBytes, bp...)
	return addrBytes
}
