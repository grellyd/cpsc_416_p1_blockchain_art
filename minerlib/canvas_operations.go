package minerlib

import (
	"blockartlib"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Operation = blockartlib.Operation
type CanvasSettings = blockartlib.CanvasSettings

const TRANSPARENT = "transparent"

// Each miner has a local instance of CanvasData
type CanvasData struct {
	Shapes []Shape
}

type Point struct {
	X, Y float64
}

type LineSegment struct {
	Point1 Point
	Point2 Point
}

func (ls *LineSegment) Length() float64 {
	return (math.Sqrt(math.Pow(ls.Point1.X-ls.Point2.X, 2) +
		math.Pow(ls.Point1.Y-ls.Point2.Y, 2)))
}

type Shape struct {
	Owner  string // Public key of owner artnode
	Hash   string
	Sides  []LineSegment
	Fill   string
	Stroke string
	Center Point
	Radius float64
}

/* Attempt to draw all shapes in list. Successfully drawn operations are in
validOps. They are attempted in a greedy fashion from the start.
validOps and invalidOps are maps. Key = shapehash, value = Operation
Handles NOP blocks. They all get added to validOps.
*/
func DrawOperations(ops []Operation, canvasSettings CanvasSettings) (validOps, invalidOps map[string]Operation, err error) {
	validOps = make(map[string]Operation)
	invalidOps = make(map[string]Operation)
	var drawnShapes []Shape
	for _, op := range ops {
		if op.Type == blockartlib.NOP {
			validOps[op.ShapeHash] = op
			continue
		}
		if op.Type == blockartlib.DELETE {
			drawnShapes, err = RemoveDrawnShape(op, drawnShapes)
			if err != nil {
				invalidOps[op.ShapeHash] = op
				return validOps, invalidOps, err
			}
			validOps[op.ShapeHash] = op
			continue
		}
		shape, err := OperationToShape(op, canvasSettings)
		if err != nil {
			return validOps, invalidOps, err
		}
		overlapsSomething := false
		for _, valid := range drawnShapes {
			if strings.Compare(valid.Owner, shape.Owner) == 0 {
				continue
			}
			if IsShapesOverlapping(valid, shape) {
				overlapsSomething = true
			}
		}
		if overlapsSomething {
			invalidOps[op.ShapeHash] = op
		} else {
			validOps[op.ShapeHash] = op
			drawnShapes = append(drawnShapes, shape)
		}
	}
	return validOps, invalidOps, err
}

// How much ink a shape needs
func InkNeeded(op Operation, settings CanvasSettings) (inkUnits uint32, err error) {
	var temp float64
	shape, err := OperationToShape(op, settings)
	if err != nil {
		return 0, err
	}
	if shape.IsCircle() {
		if shape.Fill == TRANSPARENT {
			temp = 2 * math.Pi * shape.Radius
		} else {
			temp = math.Pi * math.Pow(shape.Radius, 2)
		}
	} else {
		if shape.Fill == TRANSPARENT {
			for _, side := range shape.Sides {
				temp += side.Length()
			}
		} else {
			for i := 0; i < len(shape.Sides); i++ {
				x1 := shape.Sides[i].Point1.X
				y1 := shape.Sides[i].Point1.Y
				x2 := shape.Sides[i].Point2.X
				y2 := shape.Sides[i].Point2.Y
				temp += ((x1 * y2) - (y1 * x2))
			}
			temp /= 2
		}
	}
	inkUnits = uint32(math.Ceil(temp))
	return inkUnits, nil
}

func RemoveDrawnShape(op Operation, drawnShapes []Shape) ([]Shape, error) {
	for i, _ := range drawnShapes {
		if drawnShapes[i].Hash == op.ShapeHash {
			drawnShapes[i] = drawnShapes[len(drawnShapes)-1] 
			drawnShapes = drawnShapes[:len(drawnShapes)-1]
			return drawnShapes, nil
		}
	}
	return drawnShapes, blockartlib.InvalidShapeHashError(op.ShapeHash)
}

func ResolvePoint(initial Point, target Point, isAbsolute bool) (p Point) {
	if isAbsolute {
		return target
	}
	return AddPoints(initial, target)
}

func OperationToShape(op Operation, canvasSettings CanvasSettings) (s Shape, err error) {
	svgString := op.ShapeSVGString
	if len(svgString) > 128 {
		return s, blockartlib.ShapeSvgStringTooLongError(svgString)
	}

	re := regexp.MustCompile("[A-Za-z]|[-0-9., ]*")
	allPieces := re.FindAllString(svgString, -1)

	s.Fill = op.Fill
	s.Stroke = op.Stroke
	s.Owner = op.ArtNodePubKey
	s.Hash = op.ShapeHash

	curPt := Point{0, 0}
	initialPt := Point{0, 0}
	for index := 0; index < len(allPieces); index++ {
		command := []rune(allPieces[index])[0]
		switch unicode.ToLower(command) {
		case 'm':
			opFloats, err := StrToFloatSlice(allPieces[index+1], 2, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := Point{opFloats[0], opFloats[1]}
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			curPt = ResolvePoint(curPt, newPt, unicode.IsUpper(command))
			index++
		case 'l':
			opFloats, err := StrToFloatSlice(allPieces[index+1], 2, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := ResolvePoint(curPt, Point{opFloats[0], opFloats[1]}, unicode.IsUpper(command))
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'h':
			opFloats, err := StrToFloatSlice(allPieces[index+1], 1, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := ResolvePoint(curPt, Point{opFloats[0], curPt.Y}, unicode.IsUpper(command))
			newPt.Y = curPt.Y
			if !InBounds(newPt, canvasSettings) {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			index++
			s.Sides = append(s.Sides, LineSegment{curPt, newPt})
			curPt = newPt
		case 'v':
			opFloats, err := StrToFloatSlice(allPieces[index+1], 1, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			newPt := ResolvePoint(curPt, Point{curPt.X, opFloats[0]}, unicode.IsUpper(command))
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
			opFloats, err := StrToFloatSlice(allPieces[index+1], 2, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			s.Center = Point{opFloats[0], opFloats[1]}
			index++
		case 'r':
			opFloats, err := StrToFloatSlice(allPieces[index+1], 1, svgString)
			if err != nil {
				return s, blockartlib.InvalidShapeSvgStringError(svgString)
			}
			s.Radius = opFloats[0]
			index++
		default:
			// Get a letter that isn't one of the defined ones
			return s, blockartlib.InvalidShapeSvgStringError(svgString)
		}
	}

	// Check the given polygon is closed
	if op.Fill != TRANSPARENT && !s.IsCircle() {
		if s.Sides[0].Point1 != s.Sides[len(s.Sides)-1].Point2 {
			return s, blockartlib.InvalidShapeSvgStringError(svgString)
		}
	}

	return s, err
}

func StrToFloatSlice(s string, expectLen int, svg string) (floatSlice []float64, err error) {
	opFloats := regexp.MustCompile(",").Split(s, -1)
	if len(opFloats) != expectLen {
		return floatSlice, blockartlib.InvalidShapeSvgStringError(svg)
	}
	for _, f := range opFloats {
		floatVal, err := strconv.ParseFloat(strings.Trim(f, " "), 64)
		if err != nil {
			return floatSlice, blockartlib.InvalidShapeSvgStringError(svg)
		}
		floatSlice = append(floatSlice, floatVal)
	}
	return floatSlice, nil
}

func AddPoints(p1, p2 Point) (p Point) {
	p.X = p1.X + p2.X
	p.Y = p1.Y + p2.Y
	return p
}

func InBounds(p Point, canvasSettings CanvasSettings) (inBounds bool) {
	return !(math.Abs(p.X) > float64(canvasSettings.CanvasXMax) ||
		math.Abs(p.Y) > float64(canvasSettings.CanvasYMax) ||
		p.X < 0 || p.Y < 0)
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

func IsShapesOverlapping(s1, s2 Shape) bool {
	if !s1.IsCircle() && !s2.IsCircle() {
		for _, s := range s1.Sides {
			for _, p := range s2.Sides {
				if IsLinesIntersecting(s, p) {
					return true
				}
			}
		}
		if (s1.Fill != TRANSPARENT) && IsPointContainedInShape(s2.Sides[0].Point1, s1) ||
			(s2.Fill != TRANSPARENT && IsPointContainedInShape(s1.Sides[0].Point1, s2)) {
			return true
		}
	} else if s1.IsCircle() && s2.IsCircle() {
		return IsCircleIntersectingCircle(s1, s2)
	} else {
		if s1.IsCircle() {
			return IsCircleIntersectingPolygon(s1, s2)
		} else {
			return IsCircleIntersectingPolygon(s2, s1)
		}
	}
	return false
}

func IsPointContainedInShape(p Point, s Shape) bool {
	segment := LineSegment{p, Point{0, p.Y}}
	var numIntersections int
	prevY := s.Sides[len(s.Sides)-1].Point1.Y
	for _, side := range s.Sides {
		if side.Point1.Y == side.Point2.Y {
			continue
		}
		if IsLinesIntersecting(side, segment) {
			if side.Point1.Y == p.Y {
				if (prevY > p.Y && side.Point2.Y < p.Y) ||
					(prevY < p.Y && side.Point2.Y > p.Y) {
					numIntersections--
				}
			}
			numIntersections++
		}
		prevY = side.Point1.Y
	}
	return numIntersections%2 != 0
}

func IsLineIntersectingCircle(seg LineSegment, circle Shape) bool {
	if IsPointInCircle(seg.Point1, circle) || IsPointInCircle(seg.Point2, circle) {
		return true
	} else {
		dlen := seg.Length()
		dx := seg.Point1.X - seg.Point2.X
		dy := seg.Point1.Y - seg.Point2.Y
		segPt1 := Point{circle.Center.X + (dy/dlen)*circle.Radius, circle.Center.Y - (dx/dlen)*circle.Radius}
		segPt2 := Point{circle.Center.X - (dy/dlen)*circle.Radius, circle.Center.Y + (dx/dlen)*circle.Radius}
		cSeg := LineSegment{segPt1, segPt2}
		return IsLinesIntersecting(seg, cSeg)
	}
}

func IsPointInCircle(p Point, circle Shape) bool {
	segment := LineSegment{p, circle.Center}
	return segment.Length() <= circle.Radius
}

func (s *Shape) ShapeToSVGPath() (svg string) {
	svg += "M" + s.Sides[0].Point1.PointToString()
	for _, side := range s.Sides {
		svg += "L" + side.Point2.PointToString()
	}
	return svg
}

func IsCircleIntersectingPolygon(c, p Shape) bool {
	if c.Fill == TRANSPARENT {
		isAllPointsInCircle := true
		for _, side := range p.Sides {
			isAllPointsInCircle = IsPointInCircle(side.Point1, c)
		}
		if isAllPointsInCircle {
			return false
		}
	}
	for _, side := range p.Sides {
		if IsLineIntersectingCircle(side, c) {
			return true
		}
	}
	// Circle and polygon don't intersect
	if p.Fill == TRANSPARENT || c.Fill == TRANSPARENT {
		return false
	} else {
		return IsPointContainedInShape(c.Center, p)
	}
}

func IsCircleIntersectingCircle(c1, c2 Shape) bool {
	centerSegment := LineSegment{c1.Center, c2.Center}
	if c2.Fill == TRANSPARENT && c2.Radius > c1.Radius {
		if IsCircleContainedInCircle(c1, c2) {
			return false
		}
	} else if c1.Fill == TRANSPARENT && c1.Radius > c2.Radius {
		if IsCircleContainedInCircle(c2, c1) {
			return false
		}
	}
	return centerSegment.Length() <= (c1.Radius + c2.Radius)
}

// Returns true is c1 is contained in c2
func IsCircleContainedInCircle(c1, c2 Shape) bool {
	centerSegment := LineSegment{c1.Center, c2.Center}
	if centerSegment.Length() + c1.Radius < c2.Radius {
		return true
	}
	return false
}

func (s *Shape) IsCircle() bool {
	return s.Center.X >= 0 && s.Center.Y >= 0 && s.Radius > 0
}

func (p *Point) PointToString() (s string) {
	s += strconv.FormatFloat(p.X, 'f', -1, 64)
	s += ","
	s += strconv.FormatFloat(p.Y, 'f', -1, 64)
	return s
}