package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/protocols/proxy/http"
	"github.com/sakeven/mika/protocols/proxy/socks5"
	"github.com/sakeven/mika/protocols/transfer/kcp"
	"github.com/sakeven/mika/protocols/transfer/obfs"
	"github.com/sakeven/mika/protocols/transfer/tcp"
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
	socks5Sever := socks5.NewTCPRelay(c, servers[0].Dailer())
	socks5Sever.Serve()
}

func handleHTTP(c protocols.Protocol) {
	httpSever := http.NewRelay(c, servers[0].HTTPDailer())
	httpSever.Serve()
}

var servers []*server

type server struct {
	cg       *mika.CryptoGenerator
	address  string
	protocol string
	obfsURI  string
}

func (s *server) HTTPDailer() func(addr []byte) (protocols.Protocol, error) {
	return func(addr []byte) (protocols.Protocol, error) {
		conn, err := s.transferFactory()
		if err != nil {
			return nil, err
		}
		return mika.DialWithRawAddrHTTP(conn, addr, s.cg.NewCrypto())
	}
}

func (s *server) Dailer() func(addr []byte) (protocols.Protocol, error) {
	return func(addr []byte) (protocols.Protocol, error) {
		conn, err := s.transferFactory()
		if err != nil {
			return nil, err
		}
		return mika.DialWithRawAddr(conn, addr, s.cg.NewCrypto())
	}
}

func (s *server) transferFactory() (protocols.Protocol, error) {
	var conn protocols.Protocol
	var err error
	if s.protocol == protocols.KCP {
		return kcp.Dial(s.address)
	}

	conn, err = tcp.Dial(s.address)
	if err != nil {
		return nil, err
	}
	if s.protocol == protocols.ObfsHTTP {
		h := obfs.NewHTTP(conn, false)
		return h, h.SetURI(s.obfsURI)
	}
	return conn, nil
}

func main() {
	utils.SetLevel(utils.DebugLevel)
	conf := utils.ParseConf()
	for _, s := range conf.Server {
		if !strings.HasPrefix(s.ObfsURI, "http") {
			s.ObfsURI = "http://" + s.ObfsURI
		}

		se := &server{
			address:  fmt.Sprintf("%s:%d", s.Address, s.Port),
			cg:       mika.NewCryptoGenerator(s.Method, s.Password),
			protocol: s.Protocol,
			obfsURI:  s.ObfsURI,
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
