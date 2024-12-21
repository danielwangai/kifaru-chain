package crypto

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func newTransaction() *Transaction {
	return &Transaction{
		Data: []byte("hello world"),
	}
}

func TestSignTransaction(t *testing.T) {
	tx := newTransaction()
	privKey := GeneratePrivateKey()
	tx.Sign(privKey)
	assert.Equal(t, tx.PublicKey, privKey.PublicKey())
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	tx := newTransaction()
	privKey := GeneratePrivateKey()
	tx.Sign(privKey)
	assert.Nil(t, tx.Verify())

	//
	privKey2 := GeneratePrivateKey()
	tx.PublicKey = privKey2.PublicKey() // change public key
	assert.NotNil(t, tx.Verify())
}
