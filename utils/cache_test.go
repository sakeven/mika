package utils

import (
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite
}

func (suite *CacheSuite) TestSet() {
	key := "key"
	value := "value"
	c := NewCache()
	c.Set(key, value)
	suite.Assertions.Equal(value, c.s[key].val)
}

func (suite *CacheSuite) TestGet() {
	key := "key"
	value := "value"
	c := NewCache()
	c.Set(key, value)
	val, err := c.Get(key)
	suite.Assertions.Nil(err)
	suite.Assertions.Equal(value, val)
}

func (suite *CacheSuite) TestSetWithTTL() {
	key := "key"
	value := "value"
	c := NewCache()
	c.SetWithTTL(key, value, 10)
	val, err := c.Get(key)
	suite.Assertions.Nil(err)
	suite.Assertions.Equal(value, val)
	time.Sleep(10 * time.Second)

	val, err = c.Get(key)
	suite.Assertions.Equal(ErrNotExist, err)
	suite.Assertions.Equal("", val)
}

func (suite *CacheSuite) TestGC() {
	key := "key"
	value := "value"
	c := NewCache()
	c.gcInterval = 5 * time.Second
	c.SetWithTTL(key, value, 5)
	val, err := c.Get(key)
	suite.Assertions.Nil(err)
	suite.Assertions.Equal(value, val)
	time.Sleep(10 * time.Second)

	_, exist := c.s[key]
	suite.Assertions.False(exist)
}
