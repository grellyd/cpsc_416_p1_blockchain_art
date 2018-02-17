package minerlib

import (
	"blockartlib"
	"crypto/ecdsa"
	"fmt"
	"keys"
	//"net/http"
)

type BCStorage struct {
	BC  *Blockchain
	BCT *BCTree
}

func NewBlockchainStorage(genBlock *Block, settings *blockartlib.MinerNetSettings) *BCStorage {
	bcNode := NewBCTreeNode(genBlock, nil, 0, settings)
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
func (bcs *BCStorage) AppendBlock(block *Block, settings *blockartlib.MinerNetSettings, ourMinerPK string) (uint32, bool) {
	fmt.Println("[miner] Appending a new block to the tree: ", block, "\n")
	var ourInkToDraw uint32 = 0
	// finds the child to which to append
	// checks against blockchain if the block should go there
	//leaves := bcs.BCT.Leaves
	parentNode := FindBCTreeNode(bcs.BCT.GenesisNode, block.ParentHash)
	if parentNode == nil {
		err := fmt.Errorf("[miner]#AppendBlock: No such parent node exists!\n")
		fmt.Printf("%v\n", err)
		return ourInkToDraw, false
	}
	// fmt.Println("Parent hash ", parentNode.CurrentHash)
	d := parentNode.Depth + 1
	// fmt.Printf("k: %v\n", block.MinerPublicKey)
	k := keyToString(block.MinerPublicKey)
	// fmt.Printf("k: %v\n", k)

	var inkOnNode uint32 = parentNode.OwnerInkLvl[k]
	bcTreeNode := NewBCTreeNode(block, parentNode, inkOnNode, settings)
	bcTreeNode.Depth = d
	parentNode.Children = append(parentNode.Children, bcTreeNode)
	// TODO: add here update the block length
	// fmt.Println("BCTree after append: ", bcs.BCT.GenesisNode)

	if bcs.BC.LastNode.Current.CurrentHash == block.ParentHash {
		bccNode := NewBlockchainNode(bcTreeNode)
		//bcs.BC.LastNode.Next = bcTreeNode
		bcs.BC.LastNode.Next = bccNode
		bcs.BC.LastNode = bccNode
		// fmt.Println("BC after append: ", bcs.BC)
		PrintBC(bcs)
		return ourInkToDraw, false
	} else {
		if bcTreeNode.Depth > bcs.BC.LastNode.Current.Depth {
			nodesToInclude := walkUpToRoot(bcs.BCT, bcTreeNode)
			ourInkToDraw = rebuildBlockchain(bcs, nodesToInclude, ourMinerPK, settings)
			// fmt.Println("BC after append: ", bcs.BC)
			PrintBC(bcs)
			return ourInkToDraw, true
		}
		return ourInkToDraw, false
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

func (bcs *BCStorage) FindBlockByHash(hash string) (b *Block, err error) {
	treeNode := FindBCTreeNode(bcs.BCT.GenesisNode, hash)
	if treeNode != nil {
		return treeNode.BlockResiding, nil
	}
	return nil, nil
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

func (bcs *BCStorage) FindBlocksInBC(hashOfBlock string) []*Block {
	genNode := bcs.BC.GenesisNode
	//treeBl := genNode.Current
	blockArr := new([]*Block)
	for genNode != nil {
		if hashOfBlock == genNode.Current.CurrentHash {
			*blockArr = append(*blockArr, genNode.Current.BlockResiding)
			genNode = genNode.Next
		}

	}
	return *blockArr
}

// walks the blockchain.go Blockchain and collects all operations
func (bcs *BCStorage) Operations() (existingOps []Operation, err error) {
	currentNode := bcs.BC.GenesisNode
	for currentNode.Next != nil {
		for _, op := range currentNode.Next.Current.BlockResiding.Operations {
			existingOps = append(existingOps, *op)
		}
		currentNode = currentNode.Next
	}
	return existingOps, nil
}

// HELPER FUNCTIONS
func keyToString(key *ecdsa.PublicKey) string {
	return keys.EncodePublicKey(key)
}

func rebuildBlockchain(bcs *BCStorage, newNodes []*BCTreeNode, ourMinerHash string, settings *blockartlib.MinerNetSettings) uint32 {
	bc := bcs.BC
	var deleteInk uint32 = 0
	toRemove := walkUpToRoot(bcs.BCT, bc.LastNode.Current)
	for _, v := range toRemove {
		pk := keys.EncodePublicKey(v.BlockResiding.MinerPublicKey)
		if pk == ourMinerHash {
			if v.BlockResiding.ParentHash == bcs.BC.GenesisNode.Current.BlockResiding.ParentHash {
				continue
			}
			if len(v.BlockResiding.Operations) == 0 {
				deleteInk += settings.InkPerOpBlock
			} else {
				deleteInk += settings.InkPerNoOpBlock
			}

		}
		if len(v.BlockResiding.Operations) == 0 {
			v.OwnerInkLvl[pk] -= settings.InkPerOpBlock
		} else {
			v.OwnerInkLvl[pk] -= settings.InkPerNoOpBlock
		}

	}
	bc.LastNode = bc.GenesisNode
	for len(newNodes) != 0 {
		nn := newNodes[len(newNodes)-1]
		bcc := NewBlockchainNode(nn)
		appendBlockToBC(bc, bcc)
		newNodes = newNodes[:len(newNodes)-1]
	}
	return deleteInk
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
	bc.LastNode.Next = bccNod // updates Next for the last node
	bc.LastNode = bccNod

	return
}

/*func printTree (bct *BCTree, bctNode *BCTreeNode) {
	if bct.GenesisNode.CurrentHash == bctNode.CurrentHash
	// fmt.Println("----- Printing tree ------")
	if bct != nil {
		if len(bct.GenesisNode.Children) == 0 {
			return
		}
			for _, v := range bct.Children {
				res := FindBCTreeNode(v, nodeHash)
				if res != nil {
					return res
				}
			}
		}
	}

}*/

func PrintBC(bcs *BCStorage) {
	fmt.Println("[miner] ---Current BlockChain---")
	genNode := bcs.BC.GenesisNode

	for genNode != nil {
		// fmt.Println(genNode.Current.CurrentHash)
		genNode = genNode.Next

	}
	fmt.Println("--------------------------\n\n")
}
