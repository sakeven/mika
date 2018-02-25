package obfs

import (
	"bytes"
	"fmt"
	"net"
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
func NewHTTP(conn protocols.Protocol, uri string, isServerSide bool) *HTTP {
	return &HTTP{
		uri:          "/",
		host:         uri,
		conn:         conn,
		isServerSide: isServerSide,
	}
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
		buf := make([]byte, 0, len(b))
		for {
			n, err = h.conn.Read(b)
			if err != nil {
				return 0, err
			}
			i := bytes.Index(append(buf, b...), []byte("\r\n\r\n"))
			if i >= 0 {
				i = i - len(buf)
				n = n - i - 4
				copy(b, b[i+4:])
				break
			}
			copy(buf, b)
		}
		h.readStart = true
		return
	}

	return h.conn.Read(b)
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
