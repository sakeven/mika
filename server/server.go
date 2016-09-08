package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sakeven/ssng/ss"
	"github.com/sakeven/ssng/utils"
)

var conf *utils.Conf

func handle(c net.Conn, cg *ss.CryptoGenerate) {
	ssConn := ss.NewConn(c, cg.NewCrypto())
	ss.Serve(ssConn)
}

func Listen(serverInfo *utils.ServerConf) {
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listen on %d\n", serverInfo.Port)
	cg := ss.NewCryptoGenerate(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			handle(c, cg)
		}()
	}
}

func main() {
	conf = utils.ParseSeverConf()

	//TODO check conf

	for _, serverInfo := range conf.Server {
		Listen(serverInfo)
	}
}
