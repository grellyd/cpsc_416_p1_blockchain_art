package minerlib

type Blockchain struct {
	GenesisNode *BlockchainNode
	//NextNode	*BlockchainNode
	LastNode *BlockchainNode
}
