package crypto

type Storage interface {
	Put(b *Block) error
	// Get(b *Block) error
}

type BlockchainStore struct{}

func NewBlockchainStorage() *BlockchainStore {
	return &BlockchainStore{}
}

func (bs *BlockchainStore) Put(b *Block) error {
	return nil
}
