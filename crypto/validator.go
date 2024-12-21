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
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
