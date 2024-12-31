package network

import (
	"sort"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/types"
)

type TxMapSorter struct {
	transactions []*crypto.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*crypto.Transaction) *TxMapSorter {
	txx := make([]*crypto.Transaction, len(txMap))
	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}
	s := &TxMapSorter{txx}
	sort.Sort(s)
	return s
}
func (s *TxMapSorter) Len() int { return len(s.transactions) }
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].GetFirstSeen() < s.transactions[j].GetFirstSeen()
}

type TxPool struct {
	transactions map[types.Hash]*crypto.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*crypto.Transaction),
	}
}

func (p *TxPool) Transactions() []*crypto.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
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
