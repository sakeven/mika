package http

import (
	// "fmt"
	// "io"
	"bufio"
	"net/http"
	"net/http/httputil"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/utils"
)

type Relay struct {
	conn     protocols.Protocol
	cipher   *mika.Crypto
	ssServer string
	protocol string
	closed   bool
}

func NewRelay(conn protocols.Protocol, protocol string, mikaServer string, cipher *mika.Crypto) *Relay {
	return &Relay{
		conn:     conn,
		cipher:   cipher,
		protocol: protocol,
		ssServer: mikaServer,
	}
}

// Serve parse data and then send to mika server.
func (h *Relay) Serve() {

	bf := bufio.NewReader(h.conn)
	req, err := http.ReadRequest(bf)
	if err != nil {
		utils.Errorf("Read request error %s", err)
		return
	}

	// TODO Set http protocol flag
	mikaConn, err := mika.DailWithRawAddrHTTP(h.protocol, h.ssServer, utils.ToAddr(req.URL.Host), h.cipher)
	if err != nil {
		return
	}

	defer func() {
		if !h.closed {
			err := mikaConn.Close()
			utils.Errorf("Close connection error %v\n", err)
		}
	}()

	if req.Method == "CONNECT" {
		_HTTPSHandler(h.conn)
	} else {
		_HTTPHandler(mikaConn, req)
	}

	go protocols.Pipe(h.conn, mikaConn)
	protocols.Pipe(mikaConn, h.conn)
	h.closed = true
}

var HTTP200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

func _HTTPSHandler(client protocols.Protocol) {
	client.Write(HTTP200)
}

func _HTTPHandler(conn protocols.Protocol, req *http.Request) {
	utils.Infof("Sending request %v %v \n", req.Method, req.URL.Host)

	rmProxyHeaders(req)
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		utils.Fatalf("%s", err)
	}

	conn.Write(dump)
}

// rmProxyHeaders remove Hop-by-hop headers.
func rmProxyHeaders(req *http.Request) {
	req.RequestURI = ""
	req.Header.Del("Proxy-Connection")
	req.Header.Del("Connection")
	req.Header.Del("Keep-Alive")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("TE")
	req.Header.Del("Trailers")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Upgrade")
}
