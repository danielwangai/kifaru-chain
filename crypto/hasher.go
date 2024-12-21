package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/danielwangai/kifaru-block/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Block) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b); err != nil {
		panic(err)
	}

	h := sha256.Sum256(buf.Bytes())
	return types.Hash(h)
}
