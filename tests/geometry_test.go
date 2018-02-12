package testlib

import (
	"blockartlib"
	"minerlib"
	"testing"
)

type LineSegment = minerlib.LineSegment
type CanvasSettings = blockartlib.CanvasSettings
type Operation = blockartlib.Operation
type Point = minerlib.Point
type Shape = minerlib.Shape

func TestIntersectionsVertical(t *testing.T) {
	l1 := LineSegment{Point{5, 0}, Point{5, 5}}
	l2 := LineSegment{Point{10, 0}, Point{10, 5}}
	l3 := LineSegment{Point{5, 4}, Point{5, 7}}

	t1 := minerlib.IsLinesIntersecting(l1, l2) // Expect false
	if t1 {
		t.Errorf("intersect(l1, l2). Expected false, got %v\n", t1)
	}

	t2 := minerlib.IsLinesIntersecting(l1, l3) // Expect true
	if !t2 {
		t.Errorf("intersect(l1, l3). Expected false, got %v\n", t2)
	}
}

func TestIsLinesIntersecting(t *testing.T) {
	l1 := LineSegment{Point{5, 5}, Point{5, -5}}
	l2 := LineSegment{Point{5, 5}, Point{10, 0}}
	l3 := LineSegment{Point{5, 0}, Point{10, 5}}
	l4 := LineSegment{Point{10, -5}, Point{15, -5}}

	t1 := minerlib.IsLinesIntersecting(l1, l4) // Expect false
	if t1 {
		t.Errorf("intersect(l1, l2). Expected false, got %v\n", t1)
	}

	t2 := minerlib.IsLinesIntersecting(l1, l2) // Expect true
	if !t2 {
		t.Errorf("intersect(l1, l3). Expected false, got %v\n", t2)
	}

	t3 := minerlib.IsLinesIntersecting(l2, l3) // Expect true
	if !t3 {
		t.Errorf("intersect(l1, l3). Expected false, got %v\n", t2)
	}
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

func TestInkArea(t *testing.T) {
	svg := "M 8,0 V 8 L 4,4 l -4,4 v -8 h 8"
	op := Operation{4, 2, "opsig", blockartlib.PATH, "nonempty", "red", svg, "pubkey", 34}
	settings := CanvasSettings{100, 100}
	ink, _ := minerlib.InkNeeded(op, settings)
	if ink != 48 {
		t.Errorf("Expected ink to be 48 units. Instead was %v\n", ink)
	}

	transparentOp := Operation{4, 2, "opsig", blockartlib.PATH, "transparent", "red", svg, "pubkey", 34}
	transparentInk, _ := minerlib.InkNeeded(transparentOp, settings)
	if transparentInk != 36 {
		t.Errorf("Expected ink to be 36 units. Instead was %v\n", transparentInk)
	}
}

/*
type Shape struct {
	Owner    string // Public key of owner artnode
	Hash     string
	Sides    []LineSegment
	Fill		 string
	Stroke   string
}
*/
func TestShapesOverlappingConvave(t *testing.T) {
	so1 := LineSegment{Point{5, 2}, Point{6, 2}}
	so2 := LineSegment{Point{6, 2}, Point{6, 3}}
	so3 := LineSegment{Point{6, 3}, Point{5, 3}}
	so4 := LineSegment{Point{5, 3}, Point{5, 2}}
	soSides := []LineSegment{so1, so2, so3, so4}
	squareOut := Shape{"o", "h", soSides, "transparent", "stroke"}

	si1 := LineSegment{Point{8, 4}, Point{9, 4}}
	si2 := LineSegment{Point{9, 4}, Point{9, 5}}
	si3 := LineSegment{Point{9, 5}, Point{8, 5}}
	si4 := LineSegment{Point{8, 5}, Point{8, 4}}
	siSides := []LineSegment{si1, si2, si3, si4}
	squareIn := Shape{"o", "h", siSides, "transparent", "stroke"}

	isOverlap := minerlib.IsShapesOverlapping(squareIn, squareOut) // Expect false
	if isOverlap {
		t.Errorf("1) squareIn and squareOut. Expected false. Got %v\n", isOverlap)
	}

	c1 := LineSegment{Point{2, 2}, Point{4, 2}}
	c2 := LineSegment{Point{4, 2}, Point{4, 4}}
	c3 := LineSegment{Point{4, 4}, Point{7, 4}}
	c4 := LineSegment{Point{7, 4}, Point{7, 2}}
	c5 := LineSegment{Point{7, 2}, Point{10, 2}}
	c6 := LineSegment{Point{10, 2}, Point{10, 6}}
	c7 := LineSegment{Point{10, 6}, Point{2, 6}}
	c8 := LineSegment{Point{2, 6}, Point{2, 2}}
	cSides := []LineSegment{c1, c2, c3, c4, c5, c6, c7, c8}
	c := Shape{"o", "h", cSides, "solid", "stroke"}
	
	isOverlap = minerlib.IsShapesOverlapping(squareOut, c) // Expect false
	if isOverlap {
		t.Errorf("2) squareOut and c. Expected false. Got %v\n", isOverlap)
	}

	isOverlap = minerlib.IsShapesOverlapping(squareIn, c) // Expect true
	if !isOverlap {
		t.Errorf("3) squareIn and c. Expected true. Got %v\n", isOverlap)
	}

	c.Fill = minerlib.TRANSPARENT
	isOverlap = minerlib.IsShapesOverlapping(squareIn, c) // Expect false
	if isOverlap {
		t.Errorf("4) squareIn and c transparent. Expected false. Got %v\n", isOverlap)
	}
}