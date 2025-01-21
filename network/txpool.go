package network

import (
	"sync"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/types"
)

type TxPool struct {
	all     *TxMap
	pending *TxMap
	// maximum number of transactions in the pool
	maxLength int
}

func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxMap(),
		pending:   NewTxMap(),
		maxLength: maxLength,
	}
}

func (p *TxPool) Add(tx *crypto.Transaction) {
	// prune the oldest transaction that is sitting in the all pool
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(crypto.TxHasher{}))
	}

	if !p.all.Contains(tx.Hash(crypto.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending returns a slice of transactions that are in the pending pool
func (p *TxPool) Pending() []*crypto.Transaction {
	return p.pending.txs.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}
func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

func (p *TxPool) AllTxCount() int {
	return p.all.Count()
}

type TxMap struct {
	lock     sync.RWMutex
	txLookup map[types.Hash]*crypto.Transaction
	txs      *types.List[*crypto.Transaction]
}

func NewTxMap() *TxMap {
	return &TxMap{
		txLookup: make(map[types.Hash]*crypto.Transaction),
		txs:      types.NewList[*crypto.Transaction](),
	}
}

// First returns the first transaction in the pool
func (t *TxMap) First() *crypto.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()
	first := t.txs.Get(0)
	return t.txLookup[first.Hash(crypto.TxHasher{})]
}

// Get returns transaction in the pool matching hash
func (t *TxMap) Get(h types.Hash) *crypto.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.txLookup[h]
}

// Add inserts a new transaction to the pool
func (t *TxMap) Add(tx *crypto.Transaction) {
	hash := tx.Hash(crypto.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.txLookup[hash]; !ok {
		t.txLookup[hash] = tx
		t.txs.Insert(tx)
	}
}

// Remove deletes transaction matching hash from the pool
func (t *TxMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.txs.Remove(t.txLookup[h])
	delete(t.txLookup, h)
}

// Count returns the number of transactions in the txLookup map
func (t *TxMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.txLookup)
}

// Contains checks if a transaction matching hash is contained in the txLookup
func (t *TxMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	_, ok := t.txLookup[h]
	return ok
}

// Clear removes all transactions in the txLookup
func (t *TxMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.txLookup = make(map[types.Hash]*crypto.Transaction)
	t.txs.Clear()
}
