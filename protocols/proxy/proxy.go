package proxy

import (
	"github.com/sakeven/mika/protocols"
	"github.com/sakeven/mika/protocols/mika"
)

type ServerInfo struct {
	Protocol string
	Address  string
	cipher   *mika.Crypto
}

// Dialer creates a new protocols.Protocol conntion
type Dialer func(raw []byte) (protocols.Protocol, error)
