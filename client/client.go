package main

import (
	"log"
	"net"

	"github.com/sakeven/ssng/ss"
)

var cg = ss.NewCryptoGenerate("aes-128-cfb", "123456")

func tcp() {

	nl, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Fatal(err)
	}
	defer nl.Close()

	log.Println("listen 1080")
	for {
		c, err := nl.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(c)
	}

}

func handle(c net.Conn) {
	log.Printf("get connection %s", c.RemoteAddr())

	var cipher = cg.NewCrypto()
	socks5Sever := ss.NewSocks5TCPRelay(c, cipher)
	socks5Sever.Serve()
}

func main() {
	tcp()
}
