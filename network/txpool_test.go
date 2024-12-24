package network

import (
	"testing"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTxPoolAddTransaction(t *testing.T) {
	p := newTxPool(t)

	tx := crypto.NewTransaction([]byte("hello world"))
	p.Add(tx)
	assert.Equal(t, 1, p.Len())

	txHash := tx.Hash(crypto.TxHasher{})
	assert.True(t, p.Has(txHash))
}

func TestTxPoolFlushTransaction(t *testing.T) {
	p := newTxPool(t)

	// add transaction
	tx := crypto.NewTransaction([]byte("hello world"))
	p.Add(tx)
	assert.Equal(t, 1, p.Len())

	p.Flush()
	assert.Equal(t, 0, p.Len())
}

func TestTxPoolHandleTransaction(t *testing.T) {
	opts := ServerOpts{}
	s := NewServer(opts)

	p := newTxPool(t)
	s.memPool = p
	tx := crypto.NewTransaction([]byte("hello world"))

	_ = s.handleTransaction(tx)
	assert.Equal(t, 0, s.memPool.Len())

	// sign transaction
	privKey := crypto.GeneratePrivateKey()
	tx.Sign(privKey)
	_ = s.handleTransaction(tx)
	assert.Equal(t, 1, p.Len())
}

func newTxPool(t *testing.T) *TxPool {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
	return p
}
