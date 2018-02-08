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
	"fmt"
	"blockartlib"
	"crypto/ecdsa"
	"bytes"
	"encoding/gob"
	"minerlib/compute"
)

type Block struct {
  ParentHash string
  Operations []*blockartlib.Operation
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

// Marshalls the entire object
func (b *Block) MarshallBinary() ([]byte, error) {
	body, err := b.bodyBytes()
	if err != nil {
		return nil, fmt.Errorf("Error while marshalling block: %v", err)
	}
	return append(body, compute.Uint32AsByteArr(b.nonce)...), nil
}

// Unmarshalls bytes into a Block
func UnmarshallBinary(data []byte) (b *Block, err error) {
	var buff bytes.Buffer
	dec := gob.NewDecoder(&buff)
	b = &Block{}
	err = dec.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling block: %v", err)
	}
	return b, nil
}

// TODO: Check the reference to the nonce is maintained without returning the block
func (b* Block) Mine(difficulty uint8) error {
	if b.nonce != 0 {
		return fmt.Errorf("Block already mined!")
	}
	data, err := b.bodyBytes()
	if err != nil {
		return fmt.Errorf("Error while mining block: %v", err)
	}
	nonce, err := compute.DataConcurrent(data, difficulty)
	if err != nil {
		return fmt.Errorf("Error while mining block: %v", err)
	}
	b.nonce = nonce
	return nil
}

func (b *Block) GetHash() (hash string, err error) {
	if b.nonce == 0 {
		return "", fmt.Errorf("Block not yet mined!")
	}
	bytes, err := b.bodyBytes()
	return compute.MD5Hash(bytes, b.nonce), nil
}

func (b *Block) bodyBytes() (data []byte, err error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(b.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to encode ParentHash")
	}
	err = enc.Encode(b.Operations)
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to encode Operations")
	}
	err = enc.Encode(b.MinerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to encode MinerPublicKey")
	}
	return buff.Bytes(), nil
}
