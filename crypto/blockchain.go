package crypto

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	lock      sync.RWMutex
	logger    *logrus.Logger
	headers   []*Header
	store     Storage
	validator Validator
}

func NewBlockchain(log *logrus.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		logger:  log,
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

	for _, tx := range b.Transactions {
		bc.logger.Infof("msg: executing code, len: %d, hash: %v", len(tx.Data), tx.Hash(&TxHasher{}))
		vm := NewVM(tx.Data)
		if err := vm.Run(); err != nil {
			return err
		}
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

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

// Height returns the height of the latest block of the entire blockchain
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlock(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.headers = append(bc.headers, b.Header)
	bc.logger.WithFields(logrus.Fields{
		"msg":          "new block",
		"height":       b.Header.Height,
		"hash":         b.Hash(BlockHasher{}),
		"transactions": len(b.Transactions),
	}).Info("adding new block")
	return bc.store.Put(b)
}
