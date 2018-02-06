package minerlib

import (
	"strconv"
	"fmt"
	"blockartlib"
	"unicode"
)

type Operation = blockartlib.Operation

// Each miner has a local instance of CanvasData
type CanvasData struct {
	Points   [][]int                // TODO: Need to update the dimensions when creating canvas
	OpHashes map[string][]Operation // Operations that belong to each block
	OpOwners map[string][]Operation // Operations that belone to each artnode
}

type BCTreeNode struct {
	MinerInfo map[string]int // Map hash/identifier of miner to ink level
	// Might change int to a struct which contains more info
	BlockHash string // Hash of the block corresponding to this node
	Parent    *BCTreeNode
	Children  []*BCTreeNode
	Depth     int
}

type SvgCommand struct {
	Command rune
	X int
	Y int
}

/*
type Operation struct {
	Type OperationType
	OperationNumber int
	OperationSig string
	Shape ShapeType
	Fill string
	Stroke string
	ShapeSVGString string
	ArtNodePubKey string
	Nonce uint32
}
*/

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

func ParseSvgString(svg string) (arr []SvgCommand) {
	svg = svg + " " // To make iterating in helper easier
	for index, char := range svg {
		if unicode.IsSpace(char) {
			continue
		}
	
		// Implicit break at the end of each case
		switch char {
		// x and y
		case 'M':
			com, newIndex := GetInts(svg, index, 2, true, true, 'M')
			arr = append(arr, com)
			index = newIndex
		case 'm':
			com, newIndex := GetInts(svg, index, 2, true, true, 'm')
			arr = append(arr, com)
			index = newIndex
		case 'L':
			com, newIndex := GetInts(svg, index, 2, true, true, 'L')
			arr = append(arr, com)
			index = newIndex
		case 'l':
			com, newIndex := GetInts(svg, index, 2, true, true, 'l')
			arr = append(arr, com)
			index = newIndex
		// x
		case 'H':
			com, newIndex := GetInts(svg, index, 1, true, false, 'H')
			arr = append(arr, com)
			index = newIndex
		case 'h':
			com, newIndex := GetInts(svg, index, 1, true, false, 'h')
			arr = append(arr, com)
			index = newIndex
		// y
		case 'V':
			com, newIndex := GetInts(svg, index, 1, false, true, 'V')
			arr = append(arr, com)
			index = newIndex
		case 'v':
			com, newIndex := GetInts(svg, index, 1, false, true, 'v')
			arr = append(arr, com)
			index = newIndex
		// No parameters
		case 'Z':
			arr = append(arr, SvgCommand{'Z', -1, -1})
			index++
			fmt.Printf("case Z")
		case 'z':
			// Same as Z
			arr = append(arr, SvgCommand{'Z', -1, -1})
			fmt.Printf("case z")
			index++
		}
	}
	return arr
}

func GetInts(s string, index, numInts int, x, y bool, command rune) (sc SvgCommand, newIndex int) {
	switch numInts {
	case 1:
		for i := index; i < len(s) - 1; i++ {
			if !unicode.IsDigit([]rune(s)[i+1]) {
				if x {
					tempx, _ := strconv.Atoi(s[index:i+1])
					return SvgCommand{command, tempx, -1}, i + 1
				} else if y {
					tempy, _ := strconv.Atoi(s[index:i+1])
					return SvgCommand{command, -1, tempy}, i + 1
				}
			}
			continue
		}
	case 2:
		var tempx, tempy, tempindex int
		for i := index; i < len(s) - 1; i++ {
			if !unicode.IsDigit([]rune(s)[i+1]) {
				tempx, _ = strconv.Atoi(s[index:i+1])
				tempindex = i+1
				break			
			}
			continue
		}
		for j := tempindex; j < len(s) - 1; j++ {
			if !unicode.IsDigit([]rune(s)[j+1]) {
				tempy, _ = strconv.Atoi(s[tempindex:j+1])
				return SvgCommand{command, tempx, tempy}, j + 1
			}
			continue
		}
	}
}