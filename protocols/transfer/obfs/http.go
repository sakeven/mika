package obfs

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/utils"
)

// HTTP is http obfs
type HTTP struct {
	isServerSide bool
	readStart    bool
	writeStart   bool
	conn         protocols.Protocol
	uri          string
	host         string
}

// NewHTTP creates a http obfs connection.
// uri is "www.baicu.com/xcc"
func NewHTTP(conn protocols.Protocol, isServerSide bool) *HTTP {
	return &HTTP{
		conn:         conn,
		isServerSide: isServerSide,
	}
}

// SetURI is used at client side to configure uri.
func (h *HTTP) SetURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	h.uri = u.RequestURI()
	h.host = u.Host
	return nil
}

// RemoteAddr gets remote connection address.
func (h *HTTP) RemoteAddr() net.Addr {
	return h.conn.RemoteAddr()
}

const timeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func (h *HTTP) Write(b []byte) (n int, err error) {
	if !h.writeStart {
		var content string
		if h.isServerSide {
			// write http response header
			content = fmt.Sprintf(responseTemplate, server[0], time.Now().Format(timeFormat))
		} else {
			// write http request header
			content = fmt.Sprintf(reqTemplate, h.uri, h.host, userAgent[0])
		}
		utils.Debugf("content %s", content)
		h.conn.Write(append([]byte(content), b...))
		h.writeStart = true
		return
	}

	return h.conn.Write(b)
}

func (h *HTTP) Read(b []byte) (n int, err error) {
	if !h.readStart {
		last3 := make([]byte, 3)
		for {
			n, err = h.conn.Read(b)
			if err != nil {
				return 0, err
			}
			i := bytes.Index(append(last3, b[:n]...), []byte("\r\n\r\n"))
			if i >= 0 {
				i = i + 1
				n = n - i
				copy(b, b[i:])
				break
			}
			copyLast3(last3, b[:n])
		}
		h.readStart = true
		return
	}

	return h.conn.Read(b)
}

func copyLast3(last3, b []byte) {
	n := len(b)
	copy(last3, []byte{0, 0, 0})
	if n > 0 {
		last3[2] = b[n-1]
	}
	if n-1 > 0 {
		last3[1] = b[n-2]
	}
	if n-2 > 0 {
		last3[0] = b[n-3]
	}
}

// Close closes the connection
func (h *HTTP) Close() error {
	return h.conn.Close()
}

const reqTemplate = "POST %s HTTP/1.1\r\n" +
	"Host: %s\r\n" +
	"User-Agent: %s\r\n" +
	"Accept: */*\r\n" +
	"Content-Type: text/plain\r\n" +
	"Connection: keep-alive\r\n" +
	"\r\n"

const responseTemplate = "HTTP/1.1 200 OK\r\n" +
	"Server: %s\r\n" +
	"Date: %s\r\n" +
	"Accept: */*\r\n" +
	"\r\n"

var userAgent = []string{
	"curl/7.54.0",
}

var server = []string{
	"bfe/1.0.8.18",
	"nginx/1.2.3",
}
