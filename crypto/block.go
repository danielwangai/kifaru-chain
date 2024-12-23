package crypto

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"

	"github.com/danielwangai/kifaru-block/types"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

// Bytes transforms Header to a byte slice
func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	Header       *Header
	Transactions []Transaction
	Validator    *PublicKey
	Signature    *Signature
	hash         types.Hash
}

func NewBlock(h *Header, tx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: tx,
	}
}

func (b *Block) Encode(w io.Writer, e Encoder[*Block]) error {
	return e.Encode(w, b)
}
func (b *Block) Decode(r io.Reader, d Decoder[*Block]) error {
	return d.Decode(r, b)
}

// Sign uses the private key to sign a block
func (b *Block) Sign(privKey *PrivateKey) {
	sig := privKey.Sign(b.Header.Bytes())

	b.Validator = privKey.PublicKey()
	b.Signature = sig
}

// Verify checks the validity of the block header's signature
func (b *Block) Verify() error {
	if b.Signature == nil {
		return errors.New("block header has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return errors.New("block header has invalid signature")
	}

	// verify transactions
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}
