package proxy

import (
	"github.com/sakeven/mika/protocols"
)

// Dialer creates a new protocols.Protocol conntion
type Dialer func(raw []byte) (protocols.Protocol, error)
