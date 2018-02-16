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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"minerlib/compute"
)

/*
Let the Genesis Block be a special case of block:
    It is a block with:
		- a set hash
		- No operations
		- A nil miner key
		- A nonce of 0
	It can be identified by its nil public key.
	In other words, it's ParentHash is its assigned hash
*/

type Block struct {
	ParentHash     string
	Operations     []*blockartlib.Operation
	MinerPublicKey *ecdsa.PublicKey
	Nonce          uint32
}

func NewBlock(parentHash string, publicKey *ecdsa.PublicKey) (b *Block) {
	return &Block{
		ParentHash:     parentHash,
		MinerPublicKey: publicKey,
		Nonce:          0,
		Operations:     []*blockartlib.Operation{},
	}
}

func (b *Block) Mine(difficulty uint8) error {
	if b.Nonce != 0 {
		return fmt.Errorf("Block already mined!")
	}
	if b.MinerPublicKey == nil {
		return fmt.Errorf("Cannot mine the genesis block!")
	}
	data, err := b.bodyBytes()
	if err != nil {
		return fmt.Errorf("Error while mining block: %v", err)
	}
	nonce, err := compute.Data(data, difficulty)
	if err != nil {
		return fmt.Errorf("Error while mining block: %v", err)
	}
	b.Nonce = nonce
	return nil
}

func (b *Block) Hash() (hash string, err error) {
	if b.MinerPublicKey == nil {
		// Genesis Block
		return b.ParentHash, nil
	}
	if b.Nonce == 0 {
		return "", fmt.Errorf("Block not yet mined!")
	}
	bytes, err := b.bodyBytes()
	return compute.MD5Hash(bytes, b.Nonce), nil
}

func (b *Block) Valid(opDiff uint8, nopDiff uint8) (valid bool, err error) {
	// check genesis block
	if b.MinerPublicKey == nil && len(b.Operations) == 0 && b.Nonce == 0 {
		return true, nil
	}
	hash, err := b.Hash()
	if err != nil {
		return false, fmt.Errorf("Unable validate block: %v", err)
	}
	difficulty := uint8(0)
	if len(b.Operations) > 0 {
		difficulty = opDiff
		// check each op has a valid sig
		for _, op := range b.Operations {
			if op.ShapeHash != op.CalculateSig() {
				return false, nil
			}
		}
	} else {
		difficulty = nopDiff
	}
	// check nonce adds up
	if !compute.Valid(hash, difficulty) {
		return false, nil
	}
	return true, nil
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
func (b *Block) MarshallBinary() (marshalled [][]byte, err error) {
	// Guard against nil pubkeys
	if b.MinerPublicKey == nil {
		return nil, fmt.Errorf("Error: Unable to marshall nil public key")
	}
	gob.Register(&Block{})
	gob.Register(elliptic.P384())

	// Handle entire block
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(*b)
	if err != nil {
		return nil, fmt.Errorf("Error while encoding block in MarshallBinary: %v", err)
	}
	marshalled = append(marshalled, buff.Bytes())

	// Handle block's key separately
	//marshalled[1] = elliptic.Marshal(b.MinerPublicKey.Curve, b.MinerPublicKey.X, b.MinerPublicKey.Y)
	m1 := elliptic.Marshal(b.MinerPublicKey.Curve, b.MinerPublicKey.X, b.MinerPublicKey.Y)
	marshalled = append(marshalled, m1)
	return marshalled, nil
}

// Unmarshalls bytes into a Block
func UnmarshallBinary(data [][]byte) (b *Block, err error) {
	fmt.Println("In unmarshalling binary")
	fmt.Println("Data: ", data)
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	b = new(Block)

	// Handle block
	d1 := data[0]
	dec := gob.NewDecoder(bytes.NewReader(d1))
	err = dec.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling block: %v", err)
	}

	// Handle key separately
	b.MinerPublicKey.Curve = elliptic.P384()
	if data[1] == nil {
		return b, fmt.Errorf("data[1] is nil. Can't unmarshal")
	}
	b.MinerPublicKey.X, b.MinerPublicKey.Y = elliptic.Unmarshal(b.MinerPublicKey.Curve, data[1])

	return b, nil
}

func (b *Block) bodyBytes() (data []byte, err error) {
	// Guard against nil pubkeys
	if b.MinerPublicKey == nil {
		return nil, fmt.Errorf("Error: Unable to marshall nil public key")
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
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
