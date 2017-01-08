package main

import (
	"fmt"
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/protocols/proxy/http"
	"github.com/sakeven/mika/protocols/proxy/socks5"
	"github.com/sakeven/mika/utils"
)

func tcpServe(localConf *utils.LocalConf) {
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localConf.Address, localConf.Port))
	if err != nil {
		utils.Fatalf("Create server error %s", err)
	}
	defer nl.Close()

	utils.Infof("Client listen on %s://%s:%d", localConf.Protocol, localConf.Address, localConf.Port)

	var handleFunc func(c protocols.Protocol)
	switch localConf.Protocol {
	case protocols.HTTP:
		handleFunc = handleHTTP
	case protocols.SOCKS5:
		handleFunc = handleSocks5
	}

	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Local connection accept error %s", err)
			continue
		}
		utils.Infof("Get local connection from %s", c.RemoteAddr())
		go handleFunc(c)
	}

}

func handleSocks5(c protocols.Protocol) {
	socks5Sever := socks5.NewTCPRelay(c, servers[0].protocol, servers[0].address, servers[0].cg.NewCrypto())
	socks5Sever.Serve()
}

func handleHTTP(c protocols.Protocol) {
	httpSever := http.NewRelay(c, servers[0].protocol, servers[0].address, servers[0].cg.NewCrypto())
	httpSever.Serve()
}

type server struct {
	cg       *mika.CryptoGenerator
	address  string
	protocol string
}

var servers []*server

func main() {

	conf := utils.ParseSeverConf()
	for _, s := range conf.Server {
		se := &server{
			address:  fmt.Sprintf("%s:%d", s.Address, s.Port),
			cg:       mika.NewCryptoGenerator(s.Method, s.Password),
			protocol: s.Protocol,
		}
		servers = append(servers, se)
	}

	if len(servers) <= 0 {
		utils.Fatalf("Please configure server")
	}

	for _, localConf := range conf.Local {
		tcpServe(localConf)
	}
}
