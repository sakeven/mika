package kcp

import (
	"github.com/sakeven/mika/protocols"

	kcp "github.com/xtaci/kcp-go"
)

// Dial creats a new kcp conntion
func Dial(server string) (protocols.Protocol, error) {
	kcpConn, err := kcp.DialWithOptions(server, nil, 10, 3)
	if err != nil {
		return nil, err
	}

	kcpConn.SetStreamMode(true)
	kcpConn.SetNoDelay(1, 20, 2, 1)
	kcpConn.SetACKNoDelay(true)
	kcpConn.SetWindowSize(128, 1024)
	return kcpConn, nil
}
