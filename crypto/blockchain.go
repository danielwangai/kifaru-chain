package crypto

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	headers   []*Header
	store     Storage
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewBlockchainStorage(),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlock(genesis)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	// add block
	return bc.addBlock(b)
}

// HasBlock checks the existence of a block by height
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height() // less than/equal to since height begins at 0
}

// GetHeaderByHeight returns header at given height
// or error is height is higher than the blockchain height
func (bc *Blockchain) GetHeaderByHeight(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) is greater than the blockchain height", height)
	}

	return bc.headers[height], nil
}

// Height returns the height of the latest block of the entire blockchain
func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlock(b *Block) error {
	bc.headers = append(bc.headers, b.Header)
	logrus.WithFields(logrus.Fields{
		"height": b.Header.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")
	return bc.store.Put(b)
}
