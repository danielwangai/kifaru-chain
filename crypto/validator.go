package crypto

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Header.Height) {
		return fmt.Errorf("the blockchain already contains block of height: %d with hash %s", b.Header.Height, b.hash)
	}
	if b.Header.Height > v.bc.Height()+1 {
		return fmt.Errorf("the height of block with hash (%s) is too high", b.Hash(BlockHasher{}))
	}

	prevHeader, err := v.bc.GetHeaderByHeight(b.Header.Height - 1)
	if err != nil {
		return err
	}
	prevHash := BlockHasher{}.Hash(prevHeader)
	if prevHash != b.Header.PrevBlockHash {
		return fmt.Errorf("hash mismatch. expected previous hash: %s actual hash: %s", prevHeader.PrevBlockHash, prevHash)
	}
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
