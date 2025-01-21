package crypto

import (
	"bytes"

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
	assert.Equal(t, tx.From, privKey.PublicKey())
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	tx := newTransaction()
	privKey := GeneratePrivateKey()
	tx.Sign(privKey)
	assert.Nil(t, tx.Verify())

	//
	privKey2 := GeneratePrivateKey()
	tx.From = privKey2.PublicKey() // change public key
	assert.NotNil(t, tx.Verify())
}

func TestEncodeDecodeTransaction(t *testing.T) {
	tx := RandomTxWithSignature(t, []byte("hello world"))
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	// decode
	decoded := new(Transaction)
	assert.Nil(t, decoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, decoded)
}
