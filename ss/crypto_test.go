package ss

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	b := evpBytesToKey("123456", 16)
	log.Println(b, string(b))
}
