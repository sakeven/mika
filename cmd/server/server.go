package main

import (
	"fmt"
	"net"

	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/utils"
)

var conf *utils.Conf

func handle(c net.Conn, cg *mika.CryptoGenerate) {
	mikaConn, err := mika.NewMika(c, cg.NewCrypto(), nil)
	if err != nil {
		c.Close()
		mika.Errorf("Create mika connection error %s", err)
		return
	}
	mika.Serve(mikaConn)
}

func Listen(serverInfo *utils.ServerConf) {
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverInfo.Address, serverInfo.Port))
	if err != nil {
		mika.Fatalf("Create server error %s", err)
	}

	mika.Infof("Listen on %d\n", serverInfo.Port)
	cg := mika.NewCryptoGenerate(serverInfo.Method, serverInfo.Password)

	for {
		c, err := nl.Accept()
		if err != nil {
			mika.Errorf("Accept connection error %s", err)
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
