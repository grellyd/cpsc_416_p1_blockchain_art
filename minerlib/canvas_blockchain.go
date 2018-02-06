package minerlib

import (
	"blockartlib"
)

type Operation = blockartlib.Operation

// Each miner has a local instance of CanvasData
type CanvasData struct {
	Points [][]int // TODO: Need to update the dimensions when creating canvas
	OpHashes map[string][]Operation // Operations that belong to each block
	OpOwners map[string][]Operation // Operations that belone to each artnode
}

type BCTreeNode struct {
	MinerInfo map[string]int // Map hash/identifier of miner to ink level
	// Might change int to a struct which contains more info
	BlockHash string // Hash of the block corresponding to this node
	Parent   *BCTreeNode
	Children []*BCTreeNode
	Depth int
}

// Draw shapes from a given operation
func (cd *CanvasData) DrawOperation(op Operation) {

}

// Checks in greedy fashion if the list of provided operations can be drawn
func ValidOpList(ops []Operation) (validOps, invalidOps []Operation) {
	return validOps, invalidOps
}

// Update the canvas to the status of the blockchain at the given block hash
func (cd *CanvasData) UpdateToBlock(hash string) {

}

// Validate new block