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
	"crypto/elliptic"
	"bytes"
	"encoding/gob"
	"minerlib/compute"
)

type Block struct {
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

// TODO: Check the reference to the nonce is maintained without returning the block
func (b* Block) Mine(difficulty uint8) error {
	if b.Nonce != 0 {
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
	b.Nonce = nonce
	return nil
}

func (b *Block) GetHash() (hash string, err error) {
	if b.Nonce == 0 {
		return "", fmt.Errorf("Block not yet mined!")
	}
	bytes, err := b.bodyBytes()
	return compute.MD5Hash(bytes, b.Nonce), nil
}

// ==================
// Marshalling
// ==================

/*
Note when encoding: 
	"An interface value can be transmitted only if the concrete value itself is transmittable.
	"At least for now, that's equivalent to saying that interfaces holding typed nil pointers cannot be sent."
From: https://github.com/golang/go/issues/3704#issuecomment-66067672 by Rob Pike
Therefore, the marshall function will error when given any nil pointers
*/

// Marshalls the entire object
func (b *Block) MarshallBinary() ([]byte, error) {
	//body, err := b.bodyBytes()
	// Guard against nil pubkeys
	if b.MinerPublicKey == nil {
		return nil, fmt.Errorf("Error: Unable to marshall nil public key")
	}
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(*b)
	if err != nil {
		return nil, fmt.Errorf("Error while marshalling block: %v", err)
	}
	return buff.Bytes(), nil
	// return append(buff.Bytes(), compute.Uint32AsByteArr(b.Nonce)...), nil
}

// Unmarshalls bytes into a Block
func UnmarshallBinary(data []byte) (b *Block, err error) {
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	
	dec := gob.NewDecoder(bytes.NewReader(data))
	blkPtr := &Block{}
	err = dec.Decode(blkPtr)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling block: %v", err)
	}
	return blkPtr, nil
}


func (b *Block) bodyBytes() (data []byte, err error) {
	// Guard against nil pubkeys
	if b.MinerPublicKey == nil {
		return nil, fmt.Errorf("Error: Unable to marshall nil public key")
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	gob.Register(Block{})
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
