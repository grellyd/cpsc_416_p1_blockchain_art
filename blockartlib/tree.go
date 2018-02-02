package blockartlib

type BCTreeNode struct {
	CData     CanvasData
	MinerInfo map[string]int // Map hash/identifier of miner to ink level
	// Might change int to a struct which contains more info
	Parent   *BCTreeNode
	Children []*BCTreeNode
}
