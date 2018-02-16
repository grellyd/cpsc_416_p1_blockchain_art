package minerlib

type BCTree struct {
	GenesisNode *BCTreeNode
}

func FindBCTreeNode(bct *BCTreeNode, nodeHash string) *BCTreeNode {
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
