package utils

import (
	"bytes"

	"github.com/stretchr/testify/suite"
)

type AddressSuite struct {
	suite.Suite
}

func (suite *AddressSuite) TestGetAddress() {
	testcases := []struct {
		data         []byte
		expectedRaw  []byte
		expectedAddr string
		expectedErr  error
	}{
		{
			data:         []byte{AddrIPv4, 0x8, 0x7, 0x6, 0x5, 0x4, 0x38},
			expectedRaw:  []byte{AddrIPv4, 0x8, 0x7, 0x6, 0x5, 0x4, 0x38},
			expectedAddr: "8.7.6.5:1080",
			expectedErr:  nil,
		},
		{
			data:         []byte{AddrIPv6, 0x20, 0x1, 0xd, 0xb8, 0xa, 0xb, 0x12, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x4, 0x38},
			expectedRaw:  []byte{AddrIPv6, 0x20, 0x1, 0xd, 0xb8, 0xa, 0xb, 0x12, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x4, 0x38},
			expectedAddr: "[2001:db8:a0b:12f0::1]:1080",
			expectedErr:  nil,
		},
		{
			data:         []byte{AddrDomain, 0x15, 'h', 't', 't', 'p', 's', ':', '/', '/', 'w', 'w', 'w', '.', 'b', 'a', 'i', 'd', 'u', '.', 'c', 'o', 'm', 0x4, 0x38},
			expectedRaw:  []byte{AddrDomain, 0x15, 'h', 't', 't', 'p', 's', ':', '/', '/', 'w', 'w', 'w', '.', 'b', 'a', 'i', 'd', 'u', '.', 'c', 'o', 'm', 0x4, 0x38},
			expectedAddr: "[https://www.baidu.com]:1080",
			expectedErr:  nil,
		},
	}

	for _, t := range testcases {
		r := bytes.NewReader(t.data)
		raw, addr, err := GetAddress(r)
		suite.Assertions.Equal(t.expectedRaw, raw)
		suite.Assertions.Equal(t.expectedAddr, addr)
		suite.Assertions.Equal(t.expectedErr, err)
	}
}

func (suite *AddressSuite) TestToAddr() {
	testcases := []struct {
		data         string
		expectedAddr []byte
	}{
		{
			data:         "8.7.6.5:1080",
			expectedAddr: []byte{AddrIPv4, 0x8, 0x7, 0x6, 0x5, 0x4, 0x38},
		},
		{
			data:         "[2001:db8:a0b:12f0::1]:1080",
			expectedAddr: []byte{AddrIPv6, 0x20, 0x1, 0xd, 0xb8, 0xa, 0xb, 0x12, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x4, 0x38},
		},
		{
			data:         "[https://www.baidu.com]:1080",
			expectedAddr: []byte{AddrDomain, 0x15, 'h', 't', 't', 'p', 's', ':', '/', '/', 'w', 'w', 'w', '.', 'b', 'a', 'i', 'd', 'u', '.', 'c', 'o', 'm', 0x4, 0x38},
		},
	}

	for _, t := range testcases {
		addr := ToAddr(t.data)
		suite.Assertions.Equal(t.expectedAddr, addr)
	}
}
