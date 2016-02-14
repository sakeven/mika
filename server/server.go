package main

import (
	"log"
	"net"

	"github.com/sakeven/ssng/ss"
)

var cg = ss.NewCryptoGenerate("aes-128-cfb", "123456")

func handle(c net.Conn) {
	ssConn := ss.NewConn(c, cg.NewCrypto())
	ss.Serve(ssConn)
}

func tcp() {

	nl, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := nl.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			handle(c)
		}()
	}
}

func main() {
	tcp()
}
