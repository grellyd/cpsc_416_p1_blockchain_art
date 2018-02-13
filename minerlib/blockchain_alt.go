package minerlib

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"blockartlib"
)

type BlockchainAlt struct {
	GenesisNode *BCChainNode
	LastNode *BCChainNode
	}

type BCTree struct {
	GenesisNode *BCTreeNode
	//Leaves []*BCTreeNode
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
	//var leaves = make([]*BCTreeNode, 0)
	//leaves = append(leaves, bcNode)
	//bct := BCTree{bcNode, leaves}
	bct := BCTree{bcNode}
	bccNode := NewBCChainNode(bcNode)
	bc := BlockchainAlt{bccNode, bccNode}
	var forest Forest = make(Forest, 0)
	var bcs = BCStorage{
		&bc,
		&bct,
		&forest,
	}
	return &bcs
}

/// STUBS for expected behaviour

// REQUIRES: valid block (block that has parent on the tree among the leaves and valid in operations)
// EFFECTS: returns true if blockchain has been switched, false if blockchain is the same
func (bcs *BCStorage) AppendBlockToTree (block *Block, miner *Miner, hash string) bool {
	// finds the child to which to append
	// checks against blockchain if the block should go there
	//leaves := bcs.BCT.Leaves
	parentNode := FindBCTreeNode(bcs.BCT.GenesisNode, block.ParentHash)
	fmt.Println("Parent hash ", parentNode.CurrentHash)
	d := parentNode.Depth + 1
	k := keyToString(block.MinerPublicKey)
	var inkOnNode uint32 = parentNode.OwnerInkLvl[k]
	bcTreeNode := NewBCNodeAlt(block,parentNode, inkOnNode, miner, hash)
	bcTreeNode.Depth = d
	parentNode.Children = append(parentNode.Children, bcTreeNode)
	// TODO: add here update the block length
	fmt.Println("BCTree after append: ", bcs.BCT.GenesisNode)

	if bcs.BC.LastNode.Current.CurrentHash == block.ParentHash {
		bccNode := NewBCChainNode(bcTreeNode)
		bcs.BC.LastNode.Next = bcTreeNode
		bcs.BC.LastNode = bccNode
		fmt.Println("BC after append: ", bcs.BC)
		return false
	} else {
		if bcTreeNode.Depth > bcs.BC.LastNode.Current.Depth {
			nodesToInclude := walkUpToRoot (bcs.BCT, bcTreeNode)
			rebuildBlockchain (bcs.BC, nodesToInclude)
			fmt.Println("BC after append: ", bcs.BC)
			return true
		}
		return false
	}

	return true


	/*	var bcNode *BCTreeNode
	if len(leaves) == 0 {
		d := 1
		v := bcs.BCT.GenesisNode
		k := keyToString(block.MinerPublicKey)
		var inkOnNode uint32 = v.OwnerInkLvl[k]
		bcNode = NewBCNodeAlt(block,v, inkOnNode, miner, hash)
		bcNode.Depth = d
		fmt.Println("leaves before append: ", leaves)
		v.Children = append(v.Children, bcNode)
		leaves = append(leaves, bcNode)
		bcs.BCT.Leaves = leaves
		fmt.Println("leaves after append: ", bcs.BCT.Leaves)
		bccNode := NewBCChainNode(bcNode)
		bcs.BC.LastNode.Next = bcNode
		bcs.BC.LastNode = bccNode
		fmt.Println("BC after append: ", bcs.BC)
		return false
	}

	for i, v := range leaves {
		if v.CurrentHash == block.ParentHash {
			d := v.Depth + 1
			k := keyToString(block.MinerPublicKey)
			var inkOnNode uint32 = v.OwnerInkLvl[k]
			bcNode = NewBCNodeAlt(block,v, inkOnNode, miner, hash)
			bcNode.Depth = d
			// put a new node among the leaves
			fmt.Println("leaves before append: ", leaves)
			v.Children = append(v.Children, bcNode)
			leaves1 := append(leaves[:i], bcNode)
			leaves = append(leaves1, leaves[i+1:]...)
			bcs.BCT.Leaves = leaves
			fmt.Println("leaves after append: ", bcs.BCT.Leaves)
			break
		}

	}*/
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

// stub for the function which return children of the block
func (bcs *BCStorage) GetChildrenNodes (hashOfBlock string) ([]string, error) {
	children :=  make([]string, 0)
	node := FindBCTreeNode(bcs.BCT.GenesisNode, hashOfBlock)
	if node == nil {
		return children, blockartlib.InvalidBlockHashError("no such node on a tree")
	}
	for _, v := range node.Children {
		children = append(children, v.CurrentHash)
	}
	return children, nil
}

// REQUIRES: valid hash
func (bcs *BCStorage) AddToForest (blockHash string, block *Block) {
	forest := *bcs.For
	forest[blockHash] = block
	bcs.For = &forest
	return
}

// REQUIRES: valid hash
func (bcs *BCStorage) RemoveFromForest (blockHash string) {
	forest := *bcs.For
	delete(forest, blockHash)
	bcs.For = &forest
	return

}

func (bcs *BCStorage) IsInForest (blockHash string) bool {
	forest := *bcs.For
	_, ok := forest[blockHash]
	return ok
}

func keyToString (key *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(key)
	var c []byte = x509EncodedPub
	return string(c)
}

func rebuildBlockchain (bc *BlockchainAlt, newNodes []*BCTreeNode) {
	bc.LastNode = bc.GenesisNode
	for len(newNodes) !=0 {
		nn := newNodes[len(newNodes)-1]
		bcc := NewBCChainNode(nn)
		appendBlockToBC(bc, bcc)
		newNodes = newNodes[:len(newNodes)-1]
	}
	return
}

func walkUpToRoot (bcs *BCTree, bcn *BCTreeNode) []*BCTreeNode {
	newChildren := make([]*BCTreeNode, 0)
	for bcn.CurrentHash != bcs.GenesisNode.CurrentHash {
		newChildren = append(newChildren, bcn)
		bcn = bcn.Parent
	}
	return newChildren
}

func appendBlockToBC (bc *BlockchainAlt, bccNod *BCChainNode) {
	bc.LastNode.Next = bccNod.Current // updates Next for the last node
	bc.LastNode = bccNod
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



