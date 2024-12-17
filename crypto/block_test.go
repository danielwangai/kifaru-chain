package crypto

import (
	"bytes"
	"github.com/danielwangai/kifaru-block/examples"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var header1 Header = Header{
	Version:   1,
	PrevBlock: examples.RandomHash(),
	Timestamp: time.Now().UnixNano(),
	Height:    10,
	Nonce:     100000,
}

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: examples.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    10,
		Nonce:     100000,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlock_Encode_Decode(t *testing.T) {
	block := &Block{
		Header:       header1,
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, block.EncodeBinary(buf))

	blockDecode := &Block{}
	assert.Nil(t, blockDecode.DecodeBinary(buf))
	assert.Equal(t, block, blockDecode)
}

func TestBlock_Hash(t *testing.T) {
	block := &Block{
		Header:       header1,
		Transactions: nil,
	}
	// before hashing, hash value is zero
	assert.True(t, block.hash.IsZero())

	// hash block
	hash := block.Hash()

	assert.NotNil(t, hash)
	assert.False(t, block.hash.IsZero())
}
