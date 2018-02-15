package minerlib

import (
	"fmt"
	"blockartlib"
	"keys"
)

type BCTreeNode struct {
	OwnerInkLvl map[string]uint32 // Map [PubKey]to ink level
	BlockResiding *Block // Block residing at the current node
	Parent   *BCTreeNode // Previous node in the blockchain
	CurrentHash   string // Previous node in the blockchain
	Children []*BCTreeNode
	Depth int // length of BC
}

func NewBCTreeNode (block *Block, parent *BCTreeNode, ownerInkLvl uint32, settings *blockartlib.MinerNetSettings) *BCTreeNode {
	currHash, err := block.Hash()
	if err != nil {
		fmt.Printf("NewNode Error while hashing given block: %v", err)
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

func FindBCTreeNode (bct *BCTreeNode, nodeHash string) *BCTreeNode {
	if bct != nil {
		if bct.CurrentHash == nodeHash {
			return bct
		} else {
			if len(bct.Children) == 0 {
				return nil
			}
			for _, v := range bct.Children {
				res := FindBCTreeNode(v, nodeHash)
				if res != nil {
					return res
				}
			}
		}
	}
	return nil
}
