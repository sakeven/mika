package mika

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/sakeven/mika/utils"
)

// ------------------------------------------------------------------------
// | ver | cmd | reverse | protocol | protocol related | chunck id | hmac |
// ------------------------------------------------------------------------
// |  1  |  1  |    2    |    1     |      Variable    |    8      | 10   |
// ------------------------------------------------------------------------
type header struct {
	Ver             byte
	Cmd             byte
	Reverse         [2]byte
	Protocol        byte
	ProtocolRelated []byte
	ChunkID         uint64
	Hmac            []byte

	Addr string
}

const (
	version     byte = 0x01
	dataForward byte = 0x01

	tcpForward  byte = 0x01
	httpForward byte = 0x02
	udpForward  byte = 0x03
)

func newHeader(protocol byte, rawAddr []byte) *header {
	return &header{
		Ver:             version,
		Cmd:             dataForward,
		Protocol:        protocol,
		ProtocolRelated: rawAddr,
		ChunkID:         uint64(time.Now().Unix()),
	}
}

func (h *header) Bytes(iv []byte, key []byte) (hb []byte) {
	hb = append(hb, h.Ver)
	hb = append(hb, h.Cmd)
	hb = append(hb, h.Reverse[:]...)
	hb = append(hb, h.Protocol)
	hb = append(hb, h.ProtocolRelated...)

	chunkIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(chunkIDBytes, h.ChunkID)
	hb = append(hb, chunkIDBytes...)

	h.Hmac = hmacSha1(append(iv, key...), hb)
	hb = append(hb, h.Hmac[:]...)
	// Debugf("%#v chunk id %d", h, h.ChunkID)
	return
}

func getHeader(c io.Reader) (*header, error) {

	// raw := make([]byte, 260)
	raw := leakyBuf.Get()
	defer leakyBuf.Put(raw)

	header := new(header)

	io.ReadFull(c, raw[:5])

	// !!MUST record errors, and return first error after parse header.
	var errs []error
	// get version
	pos := 0
	if header.Ver = raw[pos]; header.Ver != version {
		errs = append(errs, fmt.Errorf("error mika version %d", header.Ver))
	}
	pos++

	header.Cmd = raw[pos]
	switch header.Cmd {
	case dataForward:
	default:
		errs = append(errs, fmt.Errorf("error mika cmd %d", header.Cmd))
		header.Cmd = dataForward
	}
	pos++

	// header.Reverse = raw[pos : pos+2]
	pos += 2

	var err error
	header.Protocol = raw[pos]
	switch header.Protocol {
	case tcpForward, httpForward:
		header.ProtocolRelated, header.Addr, err = utils.GetAddress(c)
		if err != nil {
			errs = append(errs, err)
		}
	// case httpForward:
	// 	header.ProtocolRelated = nil
	default:
		len := rand.Int() % 4096
		io.ReadFull(c, raw[:len])
		if err != nil {
			errs = append(errs, fmt.Errorf("error mika cmd %d", header.Cmd))
		}
	}

	io.ReadFull(c, raw[:18])
	header.ChunkID = binary.BigEndian.Uint64(raw[:8])
	header.Hmac = raw[8:18]

	if len(errs) > 0 {
		return nil, errs[0]
	}
	return header, header.Check()
}

func (h *header) Check() error {
	gap := time.Now().Unix() - int64(h.ChunkID)
	if gap < 0 {
		gap = -gap
	}

	if gap > 30 {
		return fmt.Errorf("Chunk expired")
	}

	return nil
}
