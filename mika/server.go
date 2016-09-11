package mika

import (
	"net"
)

func Serve(c *Mika) {
	defer c.Close()

	Infof("Connection from %s", c.RemoteAddr())

	address := c.header.Addr
	if ban(address) {
		return
	}

	Infof("Connect to %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		Errorf("Create connection error %s", err)
		return
	}

	go pipe(conn, c)
	pipe(c, conn)
}

// TODO
func ban(address string) bool {
	return false
}
