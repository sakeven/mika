package main

import (
	"net"

	"github.com/sakeven/ssng/ss"
)

var cg = ss.NewCryptoGenerate("aes-256-cfb", "password")

func tcp() {
	nl, err := net.Listen("tcp", ":1080")
	if err != nil {
		ss.Panicf("%s", err)
	}
	defer nl.Close()

	ss.Infof("Client listen on :1080")
	for {
		c, err := nl.Accept()
		if err != nil {
			ss.Errorf("Local connection accept error %s", err)
			continue
		}
		go handle(c)
	}

}

func handle(c net.Conn) {
	ss.Infof("Get local connection from %s", c.RemoteAddr())

	var cipher = cg.NewCrypto()
	socks5Sever := ss.NewSocks5TCPRelay(c, "localhost:8080", cipher)
	socks5Sever.Serve()
}

func main() {
	tcp()
}
