package ss

import (
	"net"
)

func Serve(c *Conn) {
	defer c.Close()

	raw, address, err := getAddress(c)
	if err != nil {
		Errorf("Get dest address error %s", err)
		return
	}

	Infof("Here %s %s", string(raw), address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		Errorf("Create connection error %s", err)
		return
	}

	go pipe(conn, c)
	pipe(c, conn)
}
