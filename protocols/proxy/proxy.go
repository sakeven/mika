package proxy

import (
	"github.com/sakeven/mika/protocols/mika"
)

type ServerInfo struct {
	Protocol string
	Address  string
	cipher   *mika.Crypto
}
