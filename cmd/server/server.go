package main

import (
	"fmt"
	"net"
	"time"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/protocols/transfer/obfs"
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

	utils.Infof("Listen on tcp://%s:%d\n", serverInfo.Address, serverInfo.Port)
	cg := mika.NewCryptoGenerator(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Accept connection error %s", err)
			continue
		}

		go func() {
			var tcpConn protocols.Protocol = &tcp.Conn{
				Conn:    c,
				Timeout: time.Duration(serverInfo.Timeout) * time.Second,
			}
			if serverInfo.Protocol == protocols.ObfsHTTP {
				utils.Debugf("use protocol %s", serverInfo.Protocol)
				tcpConn = obfs.NewHTTP(tcpConn, "www.baidu.com", true)
			}
			handle(tcpConn, cg)
		}()
	}
}

func listenKcp(serverInfo *utils.ServerConf) {
	nl, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port), nil, 10, 3)
	if err != nil {
		utils.Fatalf("Create server error %s", err)
	}

	utils.Infof("Listen on kcp://%s:%d\n", serverInfo.Address, serverInfo.Port)
	cg := mika.NewCryptoGenerator(serverInfo.Method, serverInfo.Password)

	for {
		conn, err := nl.AcceptKCP()
		if err != nil {
			utils.Errorf("Accept connection error %s", err)
			continue
		}

		go func() {
			conn.SetStreamMode(true)
			conn.SetNoDelay(1, 20, 2, 1)
			conn.SetACKNoDelay(true)
			conn.SetWindowSize(1024, 1024)
			kcpConn := &tcp.Conn{
				Conn:    conn,
				Timeout: time.Duration(serverInfo.Timeout) * time.Second,
			}
			handle(kcpConn, cg)
		}()
	}
}

func main() {
	conf = utils.ParseSeverConf()
	utils.SetLevel(utils.DebugLevel)

	//TODO check conf

	for _, serverInfo := range conf.Server {
		if serverInfo.Timeout <= 0 {
			serverInfo.Timeout = 30
		}

		if serverInfo.Protocol == protocols.KCP {
			listenKcp(serverInfo)
		} else {
			listen(serverInfo)
		}
	}
}
