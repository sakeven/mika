package ss

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"io"
)

type cryptoInfo struct {
	keyLen   int
	ivLen    int
	newBlock func(key []byte) cipher.Block
}

func newAesBlock(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return block
}

var cryptoInfoMap = map[string]*cryptoInfo{
	"aes-128-cfb": {16, 16, newAesBlock},
	"aes-192-cfb": {24, 16, newAesBlock},
	"aes-256-cfb": {32, 16, newAesBlock},
}

func md5Sum(data []byte) []byte {
	s := md5.Sum(data)
	return s[:]
}

func evpBytesToKey(password string, keyLen int) []byte {
	cnt := keyLen/md5.Size + 1
	ms := make([]byte, cnt*md5.Size)
	copy(ms, md5Sum([]byte(password)))

	data := make([]byte, md5.Size+len(password))

	for i := 1; i < cnt; i++ {
		// pos := i * md5.Size
		pos := i << 4
		// copy(data, ms[pos-md5.Size:pos])
		copy(data, ms[pos>>1:pos])
		copy(data[pos:], password)
		copy(ms[pos:], md5Sum(data))
	}

	return ms[:keyLen]
}

func NewCryptoGenerate(method string, password string) *CryptoGenerate {
	cryptoInfo := cryptoInfoMap[method]

	key := evpBytesToKey(password, cryptoInfo.keyLen)
	block := cryptoInfo.newBlock(key)

	return &CryptoGenerate{
		info:  cryptoInfo,
		block: block}
}

type CryptoGenerate struct {
	info  *cryptoInfo
	block cipher.Block
}

func (cg *CryptoGenerate) newEncStream() (cipher.Stream, []byte) {
	iv := make([]byte, cg.info.ivLen)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	return cipher.NewCFBEncrypter(cg.block, iv), iv
}

func (cg *CryptoGenerate) newDecStream(iv []byte) cipher.Stream {
	return cipher.NewCFBDecrypter(cg.block, iv)
}

func (cg *CryptoGenerate) NewCrypto() *Crypto {
	return &Crypto{CryptoGenerate: cg}
}

type Crypto struct {
	*CryptoGenerate
	enc cipher.Stream
	dec cipher.Stream
}

func (c *Crypto) Encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

func (c *Crypto) Decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}
