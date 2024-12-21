package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func blockWithSignature(t *testing.T, height uint32) *Block {
	privKey := GeneratePrivateKey()
	b := randomBlock(height)
	b.Sign(privKey)
	assert.NotNil(t, b.Signature)
	assert.Equal(t, privKey.PublicKey(), b.Validator)
	return b
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	b := blockWithSignature(t, 0)
	bc, err := NewBlockchain(b)
	assert.Nil(t, err)
	return bc
}

func TestBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc)
	assert.Equal(t, uint32(0), bc.Height())
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.Equal(t, uint32(0), bc.Height())

	// add another block
	b1 := blockWithSignature(t, 1)
	err := bc.AddBlock(b1)
	assert.Nil(t, err)
	// height increased to 1
	assert.Equal(t, uint32(1), bc.Height())

	// fails to add block at height already containing a block
	b2 := blockWithSignature(t, 1)
	err = bc.AddBlock(b2)
	assert.NotNil(t, err)
	// height does not update and remains as 1
	assert.Equal(t, uint32(1), bc.Height())
}
