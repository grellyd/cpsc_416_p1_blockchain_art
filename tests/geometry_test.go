package tests

import (
	"strings"
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

// TODO[sharon]: Test error case
func TestInkArea(t *testing.T) {
	svg := "M 8,0 V 8 L 4,4 l -4,4 v -8 h 8"
	// op := Operation{4, 2, "opsig", blockartlib.PATH, "nonempty", "red", svg, "pubkey", 34}
	op := Operation {
		Type: 4,
		OperationNumber: 2,
		OperationSig: "opsig",
		Shape: blockartlib.PATH,
		Fill: "nonempty",
		Stroke: "red",
		ShapeSVGString: svg,
		ArtNodePubKey: "pubkey",
		ValidateBlockNum: 34,
	}
	settings := CanvasSettings{100, 100}
	ink, _ := minerlib.InkNeeded(op, settings)
	if ink != 48 {
		t.Errorf("Expected ink to be 48 units. Instead was %v\n", ink)
	}

	// transparentOp := Operation{4, 2, "opsig", blockartlib.PATH, "transparent", "red", svg, "pubkey", 34
	transparentOp := Operation{
		Type: 4,
		OperationNumber:  2,
		OperationSig:  "opsig",
		Shape:  blockartlib.PATH,
		Fill:  "transparent",
		Stroke:  "red",
		ShapeSVGString:  svg,
		ArtNodePubKey:  "pubkey",
		ValidateBlockNum:  34,
	}
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
func TestShapesOverlappingConcave(t *testing.T) {
	squareOut := Square1()
	squareIn := Square2()

	isOverlap := minerlib.IsShapesOverlapping(squareIn, squareOut) // Expect false
	if isOverlap {
		t.Errorf("1) squareIn and squareOut. Expected false. Got %v\n", isOverlap)
	}

	c := ConvexPolygon()

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

func TestShapesOverlappingConvex(t *testing.T) {
	square := Square3()
	triangleOut := Triangle1()
	triangleIn := Triangle2()

	isOverlap := minerlib.IsShapesOverlapping(square, triangleIn) // Expect true
	if !isOverlap {
		t.Errorf("1) square and triangleIn. Expected true. Got %v\n", isOverlap)
	}

	triangleIn.Fill = minerlib.TRANSPARENT
	isOverlap = minerlib.IsShapesOverlapping(square, triangleIn) // Expect false
	if isOverlap {
		t.Errorf("2) square and triangleIn. Expected false. Got %v\n", isOverlap)
	}

	isOverlap = minerlib.IsShapesOverlapping(square, triangleOut) // Expect false
	if isOverlap {
		t.Errorf("3) square and triangleOut. Expected false. Got %v\n", isOverlap)
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

// TODO[sharon]: Add more tests, i.e. same vs different owners
func TestDrawAllShapes(t *testing.T) {
	convexPolygon := ConvexPolygon()
	convexPolygonOp := Operation{blockartlib.DRAW, 2, "convex_polygon", blockartlib.PATH, "filled", "red", convexPolygon.ShapeToSVGPath(), "artnode0", 34}
	squareOut := Square1()
	squareOutOp := Operation{blockartlib.DRAW, 2, "square_out", blockartlib.PATH, "transparent", "red", squareOut.ShapeToSVGPath(), "artnode1", 34}
	squareIn := Square2()
	squareInOp := Operation{blockartlib.DRAW, 2, "square_in", blockartlib.PATH, "transparent", "red", squareIn.ShapeToSVGPath(), "artnode2", 34}
	operations := []Operation{convexPolygonOp, squareOutOp, squareInOp}
	settings := CanvasSettings{1024, 1024}
	validOps, invalidOps, _ := minerlib.DrawOperations(operations, settings)
	validString := ConcatOps(validOps)
	invalidString := ConcatOps(invalidOps)
	AssertContains(t, validString, []string{"convex_polygon by artnode0", "square_out by artnode1"})
	AssertEquals(t, invalidString, "square_in by artnode2, ")
}

func TestDrawAllShapesWithOwnership(t *testing.T) {
	convexPolygon := ConvexPolygon()
	convexPolygonOp := Operation{blockartlib.DRAW, 2, "convex_polygon", blockartlib.PATH, "filled", "red", convexPolygon.ShapeToSVGPath(), "artnode0", 34}
	squareOut := Square1()
	squareOutOp := Operation{blockartlib.DRAW, 2, "square_out", blockartlib.PATH, "transparent", "red", squareOut.ShapeToSVGPath(), "artnode1", 34}
	squareIn := Square2()
	squareInOp := Operation{blockartlib.DRAW, 2, "square_in", blockartlib.PATH, "transparent", "red", squareIn.ShapeToSVGPath(), "artnode0", 34}
	operations := []Operation{convexPolygonOp, squareOutOp, squareInOp}
	settings := CanvasSettings{1024, 1024}
	validOps, invalidOps, _ := minerlib.DrawOperations(operations, settings)
	validString := ConcatOps(validOps)
	invalidString := ConcatOps(invalidOps)
	AssertContains(t, validString, []string{"convex_polygon by artnode0", "square_out by artnode1", "square_in by artnode0"})
	AssertEquals(t, invalidString, "")
}

func Square1() Shape {
	so1 := LineSegment{Point{5, 2}, Point{6, 2}}
	so2 := LineSegment{Point{6, 2}, Point{6, 3}}
	so3 := LineSegment{Point{6, 3}, Point{5, 3}}
	so4 := LineSegment{Point{5, 3}, Point{5, 2}}
	soSides := []LineSegment{so1, so2, so3, so4}
	return Shape{"owner1", "square1", soSides, "transparent", "stroke"}
}

func Square2() Shape {
	si1 := LineSegment{Point{8, 4}, Point{9, 4}}
	si2 := LineSegment{Point{9, 4}, Point{9, 5}}
	si3 := LineSegment{Point{9, 5}, Point{8, 5}}
	si4 := LineSegment{Point{8, 5}, Point{8, 4}}
	siSides := []LineSegment{si1, si2, si3, si4}
	return Shape{"owner2", "square2", siSides, "transparent", "stroke"}
}

func Square3() Shape {
	s1 := LineSegment{Point{4, 3}, Point{5, 3}}
	s2 := LineSegment{Point{5, 3}, Point{5, 4}}
	s3 := LineSegment{Point{5, 4}, Point{4, 4}}
	s4 := LineSegment{Point{4, 4}, Point{4, 3}}
	sides := []LineSegment{s1, s2, s3, s4}
	return Shape{"owner3", "square3", sides, "transparent", "stroke"}
}

func Triangle2() Shape {
	s1 := LineSegment{Point{6, 1}, Point{6, 6}}
	s2 := LineSegment{Point{6, 6}, Point{1, 4}}
	s3 := LineSegment{Point{1, 4}, Point{6, 1}}
	sides := []LineSegment{s1, s2, s3}
	return Shape{"owner4", "triangle1", sides, "filled", "stroke"}
}

func Triangle1() Shape {
	s1 := LineSegment{Point{2, 3}, Point{3, 5}}
	s2 := LineSegment{Point{3, 5}, Point{1, 5}}
	s3 := LineSegment{Point{1, 5}, Point{2, 3}}
	sides := []LineSegment{s1, s2, s3}
	return Shape{"owner5", "triangle2", sides, "filled", "stroke"}
}

func ConvexPolygon() Shape {
	c1 := LineSegment{Point{2, 2}, Point{4, 2}}
	c2 := LineSegment{Point{4, 2}, Point{4, 4}}
	c3 := LineSegment{Point{4, 4}, Point{7, 4}}
	c4 := LineSegment{Point{7, 4}, Point{7, 2}}
	c5 := LineSegment{Point{7, 2}, Point{10, 2}}
	c6 := LineSegment{Point{10, 2}, Point{10, 6}}
	c7 := LineSegment{Point{10, 6}, Point{2, 6}}
	c8 := LineSegment{Point{2, 6}, Point{2, 2}}
	cSides := []LineSegment{c1, c2, c3, c4, c5, c6, c7, c8}
	return Shape{"owner6", "concave_polygon", cSides, "solid", "stroke"}
}

func ConcatOps(ops map[string]Operation) string {
	opSigs := ""
	for _, op := range ops {
		opSigs += op.OperationSig + " by " + op.ArtNodePubKey + ", "
	}
	return opSigs
}

func AssertEquals(t *testing.T, s, expected string) {
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Got: %s\nExpected: %s\n", s, expected)
	}
}

func AssertContains(t *testing.T, testStr string, strArr []string) {
	for _, str := range strArr {
		if !strings.Contains(testStr, str) {
			t.Errorf("Expected string to contain %v, but not found.\n", str)
		}
	}
}
