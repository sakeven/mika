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

type HttpRelay struct {
	conn     protocols.Protocol
	cipher   *mika.Crypto
	ssServer string
	closed   bool
}

func NewHttpRelay(conn protocols.Protocol, mikaServer string, cipher *mika.Crypto) *HttpRelay {
	return &HttpRelay{
		conn:     conn,
		cipher:   cipher,
		ssServer: mikaServer,
	}
}

// HTTPRelay parse data and then send to mika server.
func (h *HttpRelay) Serve() {

	bf := bufio.NewReader(h.conn)
	req, err := http.ReadRequest(bf)
	if err != nil {
		utils.Errorf("Read request error %s", err)
		return
	}

	// TODO Set http protocol flag
	mikaConn, err := mika.DailWithRawAddrHttp("tcp", h.ssServer, utils.ToAddr(req.URL.Host), h.cipher)
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
		HttpsHandler(h.conn)
	} else {
		HttpHandler(mikaConn, req)
	}

	go protocols.Pipe(h.conn, mikaConn)
	protocols.Pipe(mikaConn, h.conn)
	h.closed = true
}

// In mika server we should parse http request.
func Handle(conn protocols.Protocol) {
	defer conn.Close()
}

var HTTP_200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

func HttpsHandler(client protocols.Protocol) {
	client.Write(HTTP_200)
}

func HttpHandler(conn protocols.Protocol, req *http.Request) {
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
