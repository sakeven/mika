package protocols

type Protocol interface {
	Write(b []byte) (n int, err error)
	Read(b []byte) (n int, err error)
	Close() error
}

type ProtocolStack struct {
	Proxy    Protocol
	Mika     Protocol
	Transfer Protocol
}
