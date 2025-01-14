package crypto

import (
	"testing"
	"time"

	"github.com/danielwangai/kifaru-block/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
		Nonce:         100000,
	}

	return NewBlock(header, []*Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b := NewBlock(header, []*Transaction{tx})
	dataHash, err := HashTransactions(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	b.Sign(privKey)

	return b
}

func TestSignBlock(t *testing.T) {
	b := randomBlock(0, types.Hash{})
	privKey := GeneratePrivateKey()
	b.Sign(privKey)
	assert.Equal(t, b.Validator, privKey.PublicKey())
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})
	assert.Nil(t, b.Verify())

	// alter block details
	b.Validator = GeneratePrivateKey().PublicKey() // different public key from the one that signed
	assert.NotNil(t, b.Verify())

	// alter block height
	b.Header.Height = 20
	assert.NotNil(t, b.Verify())
}
