package obfs

import (
	"github.com/stretchr/testify/suite"
)

type HTTPSuite struct {
	suite.Suite
}

func (suite *HTTPSuite) TestcopyLast3() {
	testcases := []struct {
		data     []byte
		expected []byte
	}{
		{
			data:     []byte{0x8, 0x7, 0x6, 0x5, 0x4, 0x38},
			expected: []byte{0x5, 0x4, 0x38},
		},
		{
			data:     []byte{0x20, 0x1},
			expected: []byte{0x0, 0x20, 0x1},
		},
		{
			data:     []byte{},
			expected: []byte{0x0, 0x0, 0x0},
		},
	}

	last3 := make([]byte, 3)
	for _, t := range testcases {
		copyLast3(last3, t.data)
		suite.Assertions.Equal(last3, t.expected)
	}
}
