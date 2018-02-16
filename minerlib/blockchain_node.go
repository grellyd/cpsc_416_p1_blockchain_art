package minerlib

type BlockchainNode struct {
	Current *BCTreeNode
	Next    *BlockchainNode
}

func NewBlockchainNode(current *BCTreeNode) *BlockchainNode {
	bccNode := BlockchainNode{
		current,
		nil,
	}
	return &bccNode
}
