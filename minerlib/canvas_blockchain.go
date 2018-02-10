package minerlib

import (
	"blockartlib"
	"math"
	"regexp"
	"strconv"
	"unicode"
)

type Operation = blockartlib.Operation
type CanvasSettings = blockartlib.CanvasSettings

const TRANSPARENT = "transparent"

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

func (ls *LineSegment) Length() float64 {
	return math.Sqrt(math.Pow(ls.Point1.X - ls.Point2.X, 2) -
		math.Pow(ls.Point1.Y - ls.Point2.Y, 2))
}

type Shape struct {
	Owner    string // Public key of owner artnode
	Hash     string
	Sides    []LineSegment
	Fill		 string
	Stroke   string
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
// Currently returns all operations
// TODO[sharon]: change return types to be map[OperationSig]Operation
func DrawOperations(ops []Operation, canvasSettings CanvasSettings) (validOps, invalidOps map[string]Operation) {
	/*var drawnShapes []Shape
	var indexMap map[string]int // maps hashes of shapes to their index in drawnShapes
	var curShape Shape
	for op := range ops {
		//curShape = StringToShape(op)
		if len(drawnShapes) == 0 {

		}
	}*/
	for _, op := range ops {
		validOps[op.OperationSig] = op
	}
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

	s.Fill = op.Fill
	s.Stroke = op.Stroke

	curPt := Point{0, 0}
	initialPt := Point{0, 0}
	var index int
	//var isTransparent bool // set based on fill
	for _, op := range opCommands {
		switch unicode.ToLower(op) {
		case 'm':
			newPt := Point{opFloats[index], opFloats[index+1]}
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index += 2
			curPt = ResolvePoint(curPt, newPt, unicode.IsUpper(op))
		case 'l':
			newPt := ResolvePoint(curPt, Point{opFloats[index], opFloats[index+1]}, unicode.IsUpper(op))
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index += 2
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'h':
			newPt := ResolvePoint(curPt, Point{opFloats[index], curPt.Y}, unicode.IsUpper(op))
			newPt.Y = curPt.Y
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'v':
			newPt := ResolvePoint(curPt, Point{curPt.X, opFloats[index]}, unicode.IsUpper(op))
			newPt.X = curPt.X
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'z':
			s.Sides = append(s.Sides, LineSegment{curPt, initialPt})
			curPt = initialPt
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

func InBounds(p Point, canvasSettings CanvasSettings) (inBounds bool) {
	return math.Abs(p.X) > float64(canvasSettings.CanvasXMax) ||
		math.Abs(p.Y) > float64(canvasSettings.CanvasYMax)
}

func IsLinesIntersecting(l1, l2 LineSegment) bool {
	o1 := Orientation(l1.Point1, l1.Point2, l2.Point1)
	o2 := Orientation(l1.Point1, l1.Point2, l2.Point2)
	o3 := Orientation(l1.Point1, l2.Point1, l2.Point2)
	o4 := Orientation(l1.Point2, l2.Point1, l2.Point2)

	if (o1 != o2 && o3 != o4) ||
		(o1 == 0 && OnSegment(l1.Point1, l2.Point1, l1.Point1)) ||
		(o2 == 0 && OnSegment(l1.Point1, l2.Point2, l1.Point1)) ||
		(o3 == 0 && OnSegment(l2.Point1, l1.Point1, l2.Point2)) ||
		(o4 == 0 && OnSegment(l2.Point1, l1.Point2, l2.Point2)) {
		return true
	}

	return false
}

func Orientation(p, q, r Point) (o int) {
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	if val == 0 {
		return 0
	}
	if val > 0 {
		return 1
	} else {
		return 2
	}
}

func OnSegment(p, q, r Point) (onSegment bool) {
	if q.X <= math.Max(p.X, r.X) &&
		q.X >= math.Min(p.X, r.X) &&
		q.Y <= math.Max(p.Y, r.Y) &&
		q.Y >= math.Min(p.Y, r.Y) {
		return true
	}
	return false
}

// How much ink a shape needs
func InkNeeded(op Operation, settings CanvasSettings) (inkUnits float64) {
	shape, _ := OperationToShape(op, settings)
	if shape.Fill == TRANSPARENT {
		for _, side := range shape.Sides {
			inkUnits += side.Length()
		}
	} else {
		for i := 0; i < len(shape.Sides); i++ {
			x1 := shape.Sides[i].Point1.X
			y1 := shape.Sides[i].Point1.Y
			x2 := shape.Sides[i].Point2.X
			y2 := shape.Sides[i].Point2.Y
			inkUnits += ((x1 * y2) - (y1 * x2))
		}
	}
	return inkUnits
}

func IsShapesOverlapping(s1, s2 Shape) bool {
	for _, s := range s1.Sides {
		for _, p := range s2.Sides {
			if IsLinesIntersecting(s, p) {
				return true
			}
		}
	}
	return false
}