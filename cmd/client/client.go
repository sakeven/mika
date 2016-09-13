package main

import (
	"fmt"
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/utils"
)

func tcpServe(conf *utils.Conf) {
	nl, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.LocalPort))
	if err != nil {
		utils.Panicf("Create server error %s", err)
	}
	defer nl.Close()

	utils.Infof("Client listen on :%d", conf.LocalPort)
	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Local connection accept error %s", err)
			continue
		}
		utils.Infof("Get local connection from %s", c.RemoteAddr())
		go handle(c)
	}

}

func handle(c protocols.Protocol) {
	socks5Sever := mika.NewSocks5TCPRelay(c, servers[0].address, servers[0].cg.NewCrypto())
	socks5Sever.Serve()
}

type server struct {
	cg      *mika.CryptoGenerate
	address string
}

var servers []*server

func main() {

	conf := utils.ParseSeverConf()
	for _, s := range conf.Server {
		se := &server{
			address: fmt.Sprintf("%s:%d", s.Address, s.Port),
			cg:      mika.NewCryptoGenerate(s.Method, s.Password),
		}
		servers = append(servers, se)
	}

	if len(servers) <= 0 {
		utils.Fatalf("Please configure server")
	}

	tcpServe(conf)
}
