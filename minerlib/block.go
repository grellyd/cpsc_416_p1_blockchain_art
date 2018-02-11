package minerlib

import (
	"blockartlib"
  "crypto/ecdsa"
)

type Block struct {
  // TODO: Figure what to put in blocks
  ParentHash string
  Operations []*blockartlib.Operation
  MinerPublicKey *ecdsa.PublicKey
  Nonce uint32
}

func NewGenesisBlock() (b Block, err error) {
  return b, err
}

// TODO: Check if block functions go here or in minerlib
func NewBlock() (b Block, err error) {
  return b, err
}

func (b *Block) Serialize() (blockBytes []byte){
  return blockBytes
}

func Deserialize(blockBytes []byte) (block *Block) {
  return block
}
