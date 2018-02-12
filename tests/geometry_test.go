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
type Shape struct {
	Owner    string // Public key of owner artnode
	Hash     string
	Sides    []LineSegment
	Fill		 string
	Stroke   string
}

type CanvasSettings struct {
	// Canvas dimensions
	CanvasXMax uint32
	CanvasYMax uint32
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
