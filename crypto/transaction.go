package crypto

import (
	"fmt"
	"io"

	"github.com/danielwangai/kifaru-block/types"
)

type Transaction struct {
	Data      []byte
	From      *PublicKey
	Signature *Signature

	hash types.Hash // cached tx hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}

	return tx.hash
}

// Sign signs a transaction
func (tx *Transaction) Sign(privKey *PrivateKey) {
	sig := privKey.Sign(tx.Data)

	tx.Signature = sig
	tx.From = privKey.PublicKey()
}

// Verify checks the validity of the transaction signature
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("transaction has invalid signature")
	}

	return nil
}

func (tx *Transaction) EncodeBinary(w io.Writer) error {
	return nil
}

func (tx *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
