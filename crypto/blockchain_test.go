package crypto

import (
	"testing"

	"github.com/danielwangai/kifaru-block/types"
	"github.com/stretchr/testify/assert"
)

func blockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := GeneratePrivateKey()
	b := randomBlock(height, prevBlockHash)
	b.Sign(privKey)
	assert.NotNil(t, b.Signature)
	assert.Equal(t, privKey.PublicKey(), b.Validator)
	return b
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	b := blockWithSignature(t, 0, types.Hash{})
	bc, err := NewBlockchain(b)
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeaderByHeight(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
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
	prevBlockHash := getPrevBlockHash(t, bc, 1)
	b1 := randomBlockWithSignature(t, 1, prevBlockHash)

	err := bc.AddBlock(b1)
	assert.Nil(t, err)
	// height increased to 1
	assert.Equal(t, uint32(1), bc.Height())

	// fails to add block at height already containing a block
	prevBlockHash1 := getPrevBlockHash(t, bc, b1.Header.Height)
	b2 := blockWithSignature(t, 1, prevBlockHash1)
	b2.Header.PrevBlockHash = b1.hash
	err = bc.AddBlock(b2)
	assert.NotNil(t, err)
	// height does not update and remains as 1
	assert.Equal(t, uint32(1), bc.Height())
}

func TestBlockTooHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc)

	// add block at a height higher than the next available slot
	b1 := blockWithSignature(t, 2, types.Hash{})
	assert.NotNil(t, bc.AddBlock(b1))
}
