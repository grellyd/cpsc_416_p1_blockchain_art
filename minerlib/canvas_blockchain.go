package minerlib

import (
	"blockartlib"
	"regexp"
	"strconv"
	"unicode"
)

type Operation = blockartlib.Operation
type CanvasSettings = blockartlib.CanvasSettings

// Each miner has a local instance of CanvasData
type CanvasData struct {
	Shapes []Shape
}

type BCTreeNode struct {
	MinerInfo map[string]int // Map hash/identifier of miner to ink level
	// Might change int to a struct which contains more info
	BlockHash string // Hash of the block corresponding to this node
	Parent    *BCTreeNode
	Children  []*BCTreeNode
	Depth     int
}

type Point struct {
	X, Y float64
}

type LineSegment struct {
	Point1 Point
	Point2 Point
}

type Shape struct {
	Owner    string // Public key of owner artnode
	Hash     string
	Sides    []LineSegment
	IsFilled bool
	Colour   string
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

// Draw all shapes in list
func DrawOperations(ops []Operation, canvasSettings CanvasSettings) (validOps, invalidOps []Operation) {
	/*var drawnShapes []Shape
	var indexMap map[string]int // maps hashes of shapes to their index in drawnShapes
	var curShape Shape
	for op := range ops {
		//curShape = StringToShape(op)
		if len(drawnShapes) == 0 {

		}
	}*/
	return validOps, invalidOps
}

func ResolvePoint(initial Point, target Point, isAbsolute bool) (p Point) {
	if isAbsolute {
		return target
	}
	return AddPoints(initial, target)
}

func OperationToShape(op Operation, canvasSettings CanvasSettings) (s Shape, err error) {
	svgString := op.ShapeSVGString
	// Turn all letters of svgString into a rune slice
	opCommands := []rune(regexp.MustCompile("[^a-zA-Z]").ReplaceAllString(svgString, ""))
	opFloatStrings := regexp.MustCompile("[^.0-9]+").Split(svgString, -1)

	var opFloats []float64
	for i := 1; i < len(opFloatStrings); i++ {
		floatVal, err := strconv.ParseFloat(opFloatStrings[i], 64)
		if err != nil {
			return s, blockartlib.InvalidShapeSvgStringError(svgString)
		}
		opFloats = append(opFloats, floatVal)
	}

	// TODO[sharon]: Add z
	var curPt Point
	var index int
	//var isTransparent bool // set based on fill
	for _, op := range opCommands {
		switch unicode.ToLower(op) {
		case 'm':
			x := opFloats[index]
			y := opFloats[index+1]
			if !InBounds(x, y, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := Point{x, y}
			index += 2
			curPt = ResolvePoint(curPt, newPt, unicode.IsUpper(op))
		case 'l':
			x := opFloats[index]
			y := opFloats[index+1]
			if !InBounds(x, y, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := Point{x, y}
			index += 2
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'h':
			x := opFloats[index]
			if !InBounds(x, -1, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := Point{x, curPt.Y}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'v':
			y := opFloats[index]
			if !InBounds(-1, y, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := Point{curPt.X, y}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'c':
			// TODO[sharon]: circle
		default:
			// Get a letter that isn't one of the defined ones
			return s, blockartlib.InvalidShapeSvgStringError(svgString)
		}
	}
	return s, err
}

func AddPoints(p1, p2 Point) (p Point) {
	p.X = p1.X + p2.X
	p.Y = p1.Y + p2.Y
	return p
}

func InBounds(x, y float64, canvasSettings CanvasSettings) (inBounds bool) {
	return x > float64(canvasSettings.CanvasXMax) || y > float64(canvasSettings.CanvasYMax)
}

func IsIntersecting(l1, l2 LineSegment) (isIntersecting bool) {
	return false
}
