package crypto

import (
	"github.com/danielwangai/kifaru-block/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func RandomTxWithSignature(t *testing.T, data []byte) *Transaction {
	privKey := GeneratePrivateKey()
	tx := NewTransaction(data)
	tx.Sign(privKey)
	assert.NotNil(t, tx.Signature)

	return tx
}

func RandomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := GeneratePrivateKey()
	tx := RandomTxWithSignature(t, []byte("hello world"))
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b := NewBlock(header, []*Transaction{tx})
	dataHash, err := HashTransactions(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	b.Sign(privKey)

	return b
}
