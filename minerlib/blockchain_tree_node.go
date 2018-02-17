package minerlib

import (
	"blockartlib"
	"fmt"
	"keys"
)

type BCTreeNode struct {
	OwnerInkLvl   map[string]uint32 // Map [PubKey]to ink level
	BlockResiding *Block            // Block residing at the current node
	Parent        *BCTreeNode       // Previous node in the blockchain
	CurrentHash   string            // Previous node in the blockchain
	Children      []*BCTreeNode
	Depth         int // length of BC
}

func NewBCTreeNode(block *Block, parent *BCTreeNode, ownerInkLvl uint32, settings *blockartlib.MinerNetSettings) *BCTreeNode {
	currHash, err := block.Hash()
	if err != nil {
		fmt.Printf("[miner]#Creating a new BlockTrain Tree Node errored  while hashing given block: %v\n", err)
	}
	var m = make(map[string]uint32)
	// If not genesis node
	if block.MinerPublicKey != nil {
		var currInkLvl uint32
		mappingKey := keys.EncodePublicKey(block.MinerPublicKey)
		if len(block.Operations) != 0 {
			currInkLvl = ownerInkLvl + settings.InkPerOpBlock
		} else {
			currInkLvl = ownerInkLvl + settings.InkPerNoOpBlock
		}
		m[mappingKey] = currInkLvl
	}
	var bcNode = BCTreeNode{
		m,
		block,
		parent,
		currHash,
		make([]*BCTreeNode, 0),
		0,
	}
	return &bcNode
}
