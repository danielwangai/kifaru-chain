package crypto

import (
	"fmt"
	"io"
)

type Transaction struct {
	Data      []byte
	PublicKey *PublicKey
	Signature *Signature
}

// Sign signs a transaction
func (tx *Transaction) Sign(privKey *PrivateKey) {
	sig := privKey.Sign(tx.Data)

	tx.Signature = sig
	tx.PublicKey = privKey.PublicKey()
}

// Verify checks the validity of the transaction signature
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.PublicKey, tx.Data) {
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
