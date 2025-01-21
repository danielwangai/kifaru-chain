package crypto

import (
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/danielwangai/kifaru-block/types"
	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	b := RandomBlockWithSignature(t, 0, types.Hash{})
	log := logrus.New()
	bc, err := NewBlockchain(log, b)
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
	b1 := RandomBlockWithSignature(t, 1, prevBlockHash)

	err := bc.AddBlock(b1)
	assert.Nil(t, err)
	// height increased to 1
	assert.Equal(t, uint32(1), bc.Height())

	// fails to add block at height already containing a block
	prevBlockHash1 := getPrevBlockHash(t, bc, b1.Header.Height)
	b2 := RandomBlockWithSignature(t, 1, prevBlockHash1)
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
	b1 := RandomBlockWithSignature(t, 2, types.Hash{})
	assert.NotNil(t, bc.AddBlock(b1))
}
