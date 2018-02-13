package minerlib

type Blockchain struct {
	GenesisHash string
	CurrentNode *BlockchainNode
}

type BlockchainNode struct {
	Prev *BlockchainNode
	Next *BlockchainNode
	Block *Block
}

func NewBlockchain(genesisHash string) (blkchain *Blockchain, err error) {
	return &Blockchain{
		GenesisHash: genesisHash,
		CurrentNode: nil,
	}, nil
}

func (b *Blockchain) AddBlock(blk *Block) (err error) {
	parent := b.CurrentNode
	newBlock := BlockchainNode{
		Prev: parent,
		Next: nil,
		Block: blk,
	}
	// TODO: handle the case of early exit
	if parent != nil {
		parentHash, err := parent.Block.Hash()
		newHash := blk.ParentHash
		if parentHash != newHash {
			// no longer head
			return err
		}
	}
	b.CurrentNode = &newBlock
	return nil
}
