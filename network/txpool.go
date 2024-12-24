package network

import (
	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/types"
)

type TxPool struct {
	transactions map[types.Hash]*crypto.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*crypto.Transaction),
	}
}

// Add stores a new transaction
func (p *TxPool) Add(tx *crypto.Transaction) error {
	hash := tx.Hash(crypto.TxHasher{})

	p.transactions[hash] = tx
	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

// Len returns the number of transactions in the transaction pool
func (p *TxPool) Len() int {
	return len(p.transactions)
}

// Flush deletes all transactions in the transaction pool
func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*crypto.Transaction)
}
