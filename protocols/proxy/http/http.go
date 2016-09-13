package http

import (
	// "fmt"
	// "io"
	"bufio"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
	"github.com/sakeven/mika/utils"
)

type HttpRelay struct {
	conn     net.Conn
	cipher   *mika.Crypto
	ssServer string
	closed   bool
}

func NewHttpRelay(conn net.Conn, mikaServer string, cipher *mika.Crypto) *HttpRelay {
	return &HttpRelay{
		conn:     conn,
		cipher:   cipher,
		ssServer: mikaServer,
	}
}

// HTTPRelay just send tcp data to mika server.
func (h *HttpRelay) Serve(rawAddr []byte) {
	// TODO Set http protoco flag
	mikaConn, err := mika.DailWithRawAddr("tcp", h.ssServer, rawAddr, h.cipher)
	if err != nil {
		return
	}

	defer func() {
		if !h.closed {
			err := mikaConn.Close()
			utils.Errorf("Close connection error %v\n", err)
		}
	}()

	go protocols.Pipe(h.conn, mikaConn)
	protocols.Pipe(mikaConn, h.conn)
	h.closed = true
}

// In mika server we should parse http request.
func (h *HttpRelay) Handle(conn protocols.Protocol) {
	// TODO Set http protoco flag
	defer conn.Close()

	bf := bufio.NewReader(conn)
	req, err := http.ReadRequest(bf)
	if err != nil {
		utils.Errorf("Read request error %s", err)
		return
	}

	if req.Method == "CONNECT" {
		HttpsHandler(conn, req)
		return
	}

	HttpHandler(conn, req)
	// h.closed = true
}

var HTTP_200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

func HttpsHandler(conn protocols.Protocol, req *http.Request) {

	remote, err := net.Dial("tcp", req.URL.Host) //建立服务端和代理服务器的tcp连接
	if err != nil {
		utils.Errorf("Failed to connect %v\n", req.RequestURI)
		// http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}

	conn.Write(HTTP_200)

	protocols.Pipe(conn, remote)
	go protocols.Pipe(remote, conn)

}

func HttpHandler(conn protocols.Protocol, req *http.Request) {
	utils.Infof("Sending request %v %v \n", req.Method, req.URL.Host)

	rmProxyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Errorf("%v", err)
		// http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
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
