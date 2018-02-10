package minerlib

import (
	//"blockartlib"
	//"fmt"
)

type BlockchainAlt struct {
	GenesisNode *BCTreeNode
	LastNode *BCTreeNode
	}

type BCTree struct {
	GenesisNode *BCTreeNode
	Leaves []*BCTreeNode
}

// maps hashes to blocks for the invalid blocks
type Forest map[string]*Block

type BCStorage struct {
	BC *BlockchainAlt
	BCT *BCTree
	For *Forest
}

func NewBCRepresentation (genBlock *Block, miner *Miner, hash string) (*BCStorage) {
	bcNode := NewBCNodeAlt(genBlock,nil, 0, miner, hash)
	var leaves = make([]*BCTreeNode, 0)
	leaves = append(leaves, bcNode)
	bct := BCTree{bcNode, leaves}
	bc := BlockchainAlt{bcNode, bcNode}
	var forest Forest = make(Forest, 0)
	var bcs = BCStorage{
		&bc,
		&bct,
		&forest,
	}
	return &bcs
}

/// STUBS for expected behaviour

// REQUIRES: valid block
// EFFECTS: returns true if blockchain has been switched, false if blockchain is the same
func (bcs *BCStorage) AppendBlockToTree (block *Block, miner *Miner, hash string) bool {
	return true
}

// stub for the function which return children of the block
func (bcs *BCStorage) GetChildrenNodes (hashOfBlock string) [] string {
	return make([]string, 0)
}

func (bcs *BCStorage) AddToForest (hash string, block *Block) {
	return
}

func (bcs *BCStorage) RemoveFromForest (hash string, block *Block) {
	return
}














/* these are invalid functions don't look at them

// REQUIRES: validation must be done prior appending
// current function only appends to the blockchain already valid block
func (bcs *BCStorage) AppendOwnBlockToTree (block *Block, miner *Miner, hash string) bool {
	//creates new node
	//appends the block to one of the leaves
	//updates the blockchain
	// in this case check for forest isn't necessary,
	// since forest is considered to be 'invalid'
	bcNode := NewBCNode(block, bc.LastNode, miner.InkLevel, miner, hash)
	appendBlock(bc, bcNode)
	fmt.Printf("node %v \n", bc.LastNode.CurrentHash)
	return true
}

// make new treenode
// add to tree

// REQUIRES: validation must be done prior appending
// current function only appends to the blockchain already valid block
func (bc *Blockchain) AppendNeighboursBlockToTree (block *Block, miner *Miner, hash string) bool  {
	//creates new node
	//checks if appending goes to forest: 'should be on a tree':
	//a) separate block, missing middle blocks -> goes to forest
	//b) missing block to some forest blocks -> upd forest (maybe merging some blocks)
	//c) block that will link forest blocks w smth on tree -->
	// --> append block to tree, append forest blocks to tree, remove blocks from forest,
	// update the BC
	// IF not in the forest:
	//appends the block to one of the leaves
	//updates the blockchain
	neghbour := findNeighbour(block, miner)
	if neghbour == nil {
		return false
	}
	bcNode := NewBCNode(block, bc.LastNode, neghbour.InkLvl, miner, hash)
	appendBlock(bc, bcNode)
	return true
}*/



