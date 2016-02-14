package ss

import (
	"log"
	"net"
)

func Serve(c *Conn) {
	defer c.Close()

	_, address, err := getAddress(c)
	if err != nil {
		log.Println(err)
		return
	}

	// log.Println(raw, address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}

	go pipe(conn, c)
	pipe(c, conn)
}
