package ss

import (
	"log"
	"net"
)

func Serve(c *Conn) {
	defer c.Close()

	raw, address, err := getAddress(c)
	if err != nil {
		log.Printf("Get dest address error %s", err)
		return
	}

	log.Println(string(raw), address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}

	go pipe(conn, c)
	pipe(c, conn)
}
