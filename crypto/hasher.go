package crypto

import (
	"crypto/sha256"

	"github.com/danielwangai/kifaru-block/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(header *Header) types.Hash {
	h := sha256.Sum256(header.Bytes())
	return types.Hash(h)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
