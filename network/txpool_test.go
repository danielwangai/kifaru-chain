package network

import (
	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxPoolAdd(t *testing.T) {
	p := NewTxPool(1)
	assert.Equal(t, 0, p.PendingCount())
	assert.Equal(t, 0, p.AllTxCount())

	tx := crypto.RandomTxWithSignature(t, []byte("hello world"))
	p.Add(tx)

	assert.Equal(t, 1, p.PendingCount())
	assert.Equal(t, 1, p.AllTxCount())

	txHash := tx.Hash(&crypto.TxHasher{})
	assert.True(t, p.all.Contains(txHash))
}

func TestTxPool_Prune(t *testing.T) {
	// maximum number of transactions in pool is 1
	p := NewTxPool(1)

	tx1 := crypto.RandomTxWithSignature(t, []byte("hello world"))
	p.Add(tx1)
	assert.Equal(t, 1, p.AllTxCount())

	txHash1 := tx1.Hash(&crypto.TxHasher{})
	assert.True(t, p.all.Contains(txHash1))

	// add another transaction. expect original tx to be pruned
	tx2 := crypto.RandomTxWithSignature(t, []byte("hello world 1"))
	p.Add(tx2)

	assert.Equal(t, 1, p.AllTxCount())

	// transaction count remains the at maximum
	assert.Equal(t, 1, p.all.Count())
}

func TestTxPool_ClearPending(t *testing.T) {
	p := NewTxPool(1)

	tx1 := crypto.RandomTxWithSignature(t, []byte("hello world"))
	p.Add(tx1)
	assert.Equal(t, 1, p.PendingCount())

	// clear pending
	p.ClearPending()
	assert.Equal(t, 0, p.PendingCount())
}
