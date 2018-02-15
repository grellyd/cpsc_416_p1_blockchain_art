package minerlib

import (
	"blockartlib"
	"crypto/ecdsa"
	"fmt"
	"keys"
)

type BCStorage struct {
	BC  *Blockchain
	BCT *BCTree
}

func NewBlockchainStorage(genBlock *Block, settings *blockartlib.MinerNetSettings) *BCStorage {
	bcNode := NewBCTreeNode(genBlock, nil, 0, settings)
	//var leaves = make([]*BCTreeNode, 0)
	//leaves = append(leaves, bcNode)
	//bct := BCTree{bcNode, leaves}
	bct := BCTree{bcNode}
	bccNode := NewBlockchainNode(bcNode)
	bc := Blockchain{bccNode, bccNode}
	var bcs = BCStorage{
		&bc,
		&bct,
	}
	return &bcs
}

// TODO: Remove bool return for err. Bool is unused and swallowed. Error would be more useful.
//		 See minerlib#MineBlocks()
// REQUIRES: valid block (block that has parent on the tree among the leaves and valid in operations)
// EFFECTS: returns true if blockchain has been switched, false if blockchain is the same
func (bcs *BCStorage) AppendBlock(block *Block, settings *blockartlib.MinerNetSettings) bool {
	// finds the child to which to append
	// checks against blockchain if the block should go there
	//leaves := bcs.BCT.Leaves
	parentNode := FindBCTreeNode(bcs.BCT.GenesisNode, block.ParentHash)
	if parentNode == nil {
		err := fmt.Errorf("No such parent node exists!")
		fmt.Printf("%v\n", err)
		return false
	}
	fmt.Println("Parent hash ", parentNode.CurrentHash)
	d := parentNode.Depth + 1
	k := keyToString(block.MinerPublicKey)
	var inkOnNode uint32 = parentNode.OwnerInkLvl[k]
	bcTreeNode := NewBCTreeNode(block, parentNode, inkOnNode, settings)
	bcTreeNode.Depth = d
	parentNode.Children = append(parentNode.Children, bcTreeNode)
	// TODO: add here update the block length
	fmt.Println("BCTree after append: ", bcs.BCT.GenesisNode)

	if bcs.BC.LastNode.Current.CurrentHash == block.ParentHash {
		bccNode := NewBlockchainNode(bcTreeNode)
		bcs.BC.LastNode.Next = bcTreeNode
		bcs.BC.LastNode = bccNode
		fmt.Println("BC after append: ", bcs.BC)
		return false
	} else {
		if bcTreeNode.Depth > bcs.BC.LastNode.Current.Depth {
			nodesToInclude := walkUpToRoot(bcs.BCT, bcTreeNode)
			rebuildBlockchain(bcs.BC, nodesToInclude)
			fmt.Println("BC after append: ", bcs.BC)
			return true
		}
		return false
	}
}

func (bcs *BCStorage) BlockPresent(b *Block) (present bool, err error) {
	hash, err := b.Hash()
	if err != nil {
		return false, fmt.Errorf("error while hashing block to check: %v", err)
	}
	treeNode := FindBCTreeNode(bcs.BCT.GenesisNode, hash)
	return treeNode != nil, nil
}

func (bcs *BCStorage) LastNodeHash() (string, error) {
	hash, err := bcs.BC.LastNode.Current.BlockResiding.Hash()
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve last node hash: %v", err)
	}
	return hash, nil
}

// function which return children of the block
func (bcs *BCStorage) GetChildrenNodes(hashOfBlock string) ([]string, error) {
	children := make([]string, 0)
	node := FindBCTreeNode(bcs.BCT.GenesisNode, hashOfBlock)
	if node == nil {
		return children, blockartlib.InvalidBlockHashError("no such node on a tree")
	}
	for _, v := range node.Children {
		children = append(children, v.CurrentHash)
	}
	return children, nil
}

// HELPER FUNCTIONS
func keyToString(key *ecdsa.PublicKey) string {
	return keys.EncodePublicKey(key)
}

func rebuildBlockchain(bc *Blockchain, newNodes []*BCTreeNode) {
	bc.LastNode = bc.GenesisNode
	for len(newNodes) != 0 {
		nn := newNodes[len(newNodes)-1]
		bcc := NewBlockchainNode(nn)
		appendBlockToBC(bc, bcc)
		newNodes = newNodes[:len(newNodes)-1]
	}
	return
}

func walkUpToRoot(bcs *BCTree, bcn *BCTreeNode) []*BCTreeNode {
	newChildren := make([]*BCTreeNode, 0)
	for bcn.CurrentHash != bcs.GenesisNode.CurrentHash {
		newChildren = append(newChildren, bcn)
		bcn = bcn.Parent
	}
	return newChildren
}

func appendBlockToBC(bc *Blockchain, bccNod *BlockchainNode) {
	bc.LastNode.Next = bccNod.Current // updates Next for the last node
	bc.LastNode = bccNod
	return
}
