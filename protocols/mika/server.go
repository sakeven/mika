package mika

import (
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/utils"
)

func Serve(c *Mika) {
	defer c.Close()

	utils.Infof("Connection from %s", c.RemoteAddr())

	address := c.header.Addr
	if ban(address) {
		return
	}

	utils.Infof("Connect to %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		utils.Errorf("Create connection error %s", err)
		return
	}

	go protocols.Pipe(conn, c)
	protocols.Pipe(c, conn)
}

func TCPServe(c *Mika) {
	defer c.Close()

	utils.Infof("Connection from %s", c.RemoteAddr())

	address := c.header.Addr
	if ban(address) {
		return
	}

	utils.Infof("Connect to %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		utils.Errorf("Create connection error %s", err)
		return
	}

	go protocols.Pipe(conn, c)
	protocols.Pipe(c, conn)
}

// TODO
func ban(address string) bool {
	return false
}
