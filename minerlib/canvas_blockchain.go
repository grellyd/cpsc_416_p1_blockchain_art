package minerlib

type CanvasData struct {
	Points [][]int // TODO: Need to update the dimensions when creating canvas
	//Operations map[string]Operation <- (Graham: why do we need this?)
}

type BCTreeNode struct {
	CData     CanvasData
	MinerInfo map[string]int // Map hash/identifier of miner to ink level
	// Might change int to a struct which contains more info
	Parent   *BCTreeNode
	Children []*BCTreeNode
}

func (n *BCTreeNode) validPicture() (valid bool, err error) {
	return false, nil
}
