package network

import (
	"strconv"
	"testing"

	"math/rand"

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

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := crypto.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}
	assert.Equal(t, txLen, p.Len())
	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].GetFirstSeen() < txx[i+1].GetFirstSeen())
	}
}
