package crypto

import (
	"fmt"

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

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
