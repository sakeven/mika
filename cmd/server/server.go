package main

import (
	"fmt"
	"net"
	"time"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/protocols/transfer/tcp"
	"github.com/sakeven/mika/utils"

	"github.com/xtaci/kcp-go"
)

var conf *utils.Conf

func handle(c protocols.Protocol, cg *mika.CryptoGenerator) {
	mikaConn, err := mika.NewMika(c, cg.NewCrypto(), nil)
	if err != nil {
		c.Close()
		utils.Errorf("Create mika connection error %s", err)
		return
	}
	mika.Serve(mikaConn)
}

func listen(serverInfo *utils.ServerConf) {
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port))
	if err != nil {
		utils.Fatalf("Create server error %s", err)
	}

	utils.Infof("Listen on %d\n", serverInfo.Port)
	cg := mika.NewCryptoGenerator(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Accept connection error %s", err)
			continue
		}

		go func() {
			tcpConn := &tcp.Conn{c, time.Duration(serverInfo.Timeout) * time.Second}
			handle(tcpConn, cg)
		}()
	}
}

func listenKcp(serverInfo *utils.ServerConf) {
	nl, err := kcp.Listen(fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port))
	if err != nil {
		utils.Fatalf("Create server error %s", err)
	}

	utils.Infof("Listen on kcp://%s:%d\n", serverInfo.Address, serverInfo.Port)
	cg := mika.NewCryptoGenerator(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Accept connection error %s", err)
			continue
		}

		go func() {
			tcpConn := &tcp.Conn{c, time.Duration(serverInfo.Timeout) * time.Second}
			handle(tcpConn, cg)
		}()
	}
}

func main() {
	conf = utils.ParseSeverConf()

	//TODO check conf

	for _, serverInfo := range conf.Server {
		if serverInfo.Timeout <= 0 {
			serverInfo.Timeout = 30
		}
		listenKcp(serverInfo)
	}
}
