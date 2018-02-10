package minerlib

import (
	"crypto/ecdsa"
	"crypto/x509"
)

type BCTreeNode struct {
	OwnerInkLvl map[string]uint32 // Map [PubKey]to ink level
	// Might change int to a struct which contains more info
	//BlockHash string // Hash of the block corresponding to this node
	BlockResiding *Block // Block residing at the current node
	Parent   *BCTreeNode // Previous node in the blockchain
	CurrentHash   string // Previous node in the blockchain
	Children []*BCTreeNode
	Depth int // length of BC
}

func NewBCNodeAlt (block *Block, parent *BCTreeNode, ownerInkLvl uint32, miner *Miner, currHash string) *BCTreeNode {
	var m = make(map[string]uint32)
	var currInkLvl uint32
	mappingKey := encode1(block.MinerPublicKey)
	if len(block.Operations) != 0 {
		currInkLvl = ownerInkLvl + miner.Settings.InkPerOpBlock
	} else {
		currInkLvl = ownerInkLvl + miner.Settings.InkPerNoOpBlock
	}
	m[mappingKey] = currInkLvl
	var bcNode = BCTreeNode{
		m,
		block,
		parent,
		currHash,
		nil,
		0,
	}
	return &bcNode
}

func encode1(publicKey *ecdsa.PublicKey) (string) {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	var c []byte = x509EncodedPub
	return string(c)

}