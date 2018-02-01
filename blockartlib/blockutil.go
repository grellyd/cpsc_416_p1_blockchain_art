package blockartlib

type Block struct {
}

func NewGenesisBlock() (b block, err error) {
}

func NewBlock() (b block, err error) {
  return b, err
}

func (b *Block) Serialize() (blockBytes []byte){
  return blockBytes
}

func Deserialize(blockBytes []byte) (block *Block) {
  return block
}
