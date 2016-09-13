package main

import (
	"fmt"
	"net"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/utils"
)

var conf *utils.Conf

func handle(c protocols.Protocol, cg *mika.CryptoGenerate) {
	mikaConn, err := mika.NewMika(c, cg.NewCrypto(), nil)
	if err != nil {
		c.Close()
		utils.Errorf("Create mika connection error %s", err)
		return
	}
	mika.Serve(mikaConn)
}

func Listen(serverInfo *utils.ServerConf) {
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port))
	if err != nil {
		utils.Fatalf("Create server error %s", err)
	}

	utils.Infof("Listen on %d\n", serverInfo.Port)
	cg := mika.NewCryptoGenerate(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			utils.Errorf("Accept connection error %s", err)
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
