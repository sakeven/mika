package main

import (
	"net"

	"github.com/sakeven/mika/mika"
)

var cg = mika.NewCryptoGenerate("aes-256-cfb", "password")

func tcp() {
	nl, err := net.Listen("tcp", ":1080")
	if err != nil {
		mika.Panicf("%s", err)
	}
	defer nl.Close()

	mika.Infof("Client listen on :1080")
	for {
		c, err := nl.Accept()
		if err != nil {
			mika.Errorf("Local connection accept error %s", err)
			continue
		}
		go handle(c)
	}

}

func handle(c net.Conn) {
	mika.Infof("Get local connection from %s", c.RemoteAddr())

	var cipher = cg.NewCrypto()
	socks5Sever := mika.NewSocks5TCPRelay(c, ":8388", cipher)
	socks5Sever.Serve()
}

func main() {
	tcp()
}
