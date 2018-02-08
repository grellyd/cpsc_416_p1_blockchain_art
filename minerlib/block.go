/*
Blocks are:
	- A hash of the previous block in the chain (prev-hash)
 	- An unordered set of operation records; each operation record should include:
 		- An application shape operation (op)
 		- A signature of the operation (op-sig)
 		- A public key of the art node that generated the op (used to validate op/op-sig)
	- The public key of the miner that computed this block (pub-keyMiner)
	- A 32-bit unsigned integer nonce (nonce)

The set of operations will either be valid operations (OPs) or empty operations (NOPs) but that knowledge is unimportant for validating the blocks.
*/

package minerlib

import (
	"blockartlib"
	"crypto/ecdsa"
)

type Block struct {
  Operations []*blockartlib.Operation
  ParentHash string
  MinerPublicKey *ecdsa.PublicKey
  nonce uint32
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

func (b *Block) ToBytes() (bytes []byte) {
	return nil
}
