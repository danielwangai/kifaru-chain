package crypto

import (
	"github.com/danielwangai/kifaru-block/examples"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var header1 Header = Header{}

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:   1,
		PrevBlock: examples.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    height,
		Nonce:     100000,
	}
	tx := Transaction{
		Data: []byte("hello world"),
	}

	return NewBlock(header, []Transaction{tx})
}

func TestSignBlock(t *testing.T) {
	b := randomBlock(0)
	privKey := GeneratePrivateKey()
	b.Sign(privKey)
	assert.Equal(t, b.Validator, privKey.PublicKey())
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	b := randomBlock(0)
	privKey := GeneratePrivateKey()
	b.Sign(privKey)
	assert.Nil(t, b.Verify())

	// alter block details
	b.Validator = GeneratePrivateKey().PublicKey() // different public key from the one that signed
	assert.NotNil(t, b.Verify())

	// alter block height
	b.Header.Height = 20
	assert.NotNil(t, b.Verify())
}
