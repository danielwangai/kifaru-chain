package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"time"

	"github.com/danielwangai/kifaru-block/types"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	DataHash      types.Hash
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
	*Header
	Transactions []*Transaction
	Validator    *PublicKey
	Signature    *Signature
	hash         types.Hash
}

func NewBlock(h *Header, txs []*Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

// NewBlockFromPrevHeader ...
// TODO: determine how many transaction fit in a block
func NewBlockFromPrevHeader(prevHeader *Header, txs []*Transaction) (*Block, error) {
	dataHash, err := HashTransactions(txs)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, txs), nil
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
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

	//verify data hash
	dataHash, err := HashTransactions(b.Transactions)
	if err != nil {
		return err
	}

	if b.DataHash != dataHash {
		return errors.New("block data hash does not match")
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
	b.Transactions = append(b.Transactions, tx)
}

// HashTransactions computes the hash of transaction(s) in a block
func HashTransactions(txs []*Transaction) (types.Hash, error) {
	buf := &bytes.Buffer{}

	for _, tx := range txs {
		if err := tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return types.Hash{}, err
		}
	}

	hash := sha256.Sum256(buf.Bytes())

	return hash, nil
}
