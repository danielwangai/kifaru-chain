package crypto

import (
	"bytes"
	"github.com/danielwangai/kifaru-block/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignBlock(t *testing.T) {
	b := RandomBlockWithSignature(t, 0, types.Hash{})
	privKey := GeneratePrivateKey()
	b.Sign(privKey)
	assert.Equal(t, b.Validator, privKey.PublicKey())
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	b := RandomBlockWithSignature(t, 0, types.Hash{})
	assert.Nil(t, b.Verify())

	// alter block details
	b.Validator = GeneratePrivateKey().PublicKey() // different public key from the one that signed
	assert.NotNil(t, b.Verify())

	// alter block height
	b.Header.Height = 20
	assert.NotNil(t, b.Verify())
}

func TestBlock_EncodeDecode(t *testing.T) {
	b := RandomBlockWithSignature(t, 0, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecode)
}
