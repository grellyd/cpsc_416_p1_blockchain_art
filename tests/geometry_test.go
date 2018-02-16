package tests

import (
	"blockartlib"
	"minerlib"
	"testing"
)

func TestIntersectionsVertical(t *testing.T) {
	l1 := LineSegment{Point{5, 0}, Point{5, 5}}
	l2 := LineSegment{Point{10, 0}, Point{10, 5}}
	l3 := LineSegment{Point{5, 4}, Point{5, 7}}

	t1 := minerlib.IsLinesIntersecting(l1, l2) // Expect false
	ExpectFalse(t, t1, "intersect(l1, l2). Expected false, got true.\n")

	t2 := minerlib.IsLinesIntersecting(l1, l3) // Expect true
	ExpectTrue(t, t2, "intersect(l1, l3). Expected false, got true.\n")
}

func TestIsLinesIntersecting(t *testing.T) {
	l1 := LineSegment{Point{5, 5}, Point{5, -5}}
	l2 := LineSegment{Point{5, 5}, Point{10, 0}}
	l3 := LineSegment{Point{5, 0}, Point{10, 5}}
	l4 := LineSegment{Point{10, -5}, Point{15, -5}}

	t1 := minerlib.IsLinesIntersecting(l1, l4) // Expect false
	ExpectFalse(t, t1, "intersect(l1, l2). Expected false, got true.\n")

	t2 := minerlib.IsLinesIntersecting(l1, l2) // Expect true
	ExpectTrue(t, t2, "intersect(l1, l3). Expected false. Got true.\n")

	t3 := minerlib.IsLinesIntersecting(l2, l3) // Expect true
	ExpectTrue(t, t3, "intersect(l1, l3). Expected false. Got true.\n")
}

func TestInkArea(t *testing.T) {
	svg := "M 8,0 V 8 L 4,4 l -4,4 v -8 h 8"
	// op := Operation{4, 2, "opsig", blockartlib.PATH, "nonempty", "red", svg, "pubkey", 34}
	op := Operation {
		Type: 4,
		OperationNumber: 2,
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

func TestOpToCircleAndCircleInk(t *testing.T) {
	svg := "c 10,10 r 5"
	op := Operation{4, 2, blockartlib.PATH, "transparent", "red", svg, "pubkey", 28, "shapehash", 129}
	settings := CanvasSettings{100, 100}
	circle, _ := minerlib.OperationToShape(op, settings)
	if circle.Radius != 5 {
		t.Errorf("Circle didn't work")
	}

	// Empty circle
	ink, _ := minerlib.InkNeeded(op, settings)
	if ink != 32 {
		t.Errorf("Expected 32 ink. Got: %v\n", ink)
	}

	// Filled circle
	op.Fill = "filled"
	ink, _ = minerlib.InkNeeded(op, settings)
	if ink != 79 {
		t.Errorf("Expected 79 ink. Got: %v\n", ink)
	}
}

func TestOutOfBoundsPoint(t *testing.T) {
	// L -4,4 is out of bounds
	svg := "M 8,0 V 8 L 4,4 L -4,4 v -8 h 8"
	op := Operation{4, 2, blockartlib.PATH, "nonempty", "red", svg, "pubkey", 28, "shapehash", 129}
	settings := CanvasSettings{100, 100}
	_, err := minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Bad shape svg string [M 8,0 V 8 L 4,4 L -4,4 v -8 h 8]")

	// M 1000 out of bounds
	op.ShapeSVGString = "M 1000,0 V 8 L 4,4 l -4,4 v -8 h 8"
	_, err = minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Bad shape svg string [M 1000,0 V 8 L 4,4 l -4,4 v -8 h 8]")
}

func TestInvalidSVGString(t *testing.T) {
	// Too many numbers after M
	svg := "M 8,0,0 V 8 L 4,4 l -4,4 v -8 h 8"
	op := Operation{4, 2, blockartlib.PATH, "nonempty", "red", svg, "pubkey", 34, "shapehash", 129}
	settings := CanvasSettings{100, 100}
	_, err := minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Bad shape svg string [M 8,0,0 V 8 L 4,4 l -4,4 v -8 h 8]")

	// q is not a valid command
	op.ShapeSVGString = "M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 8"
	_, err = minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Bad shape svg string [M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 8]")

	// Polygon doesn't start and end at the same place
	op.ShapeSVGString = "M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 6"
	_, err = minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Bad shape svg string [M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 6]")

	// String too long
	op.ShapeSVGString = "M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8" // q is not a valid command
	_, err = minerlib.OperationToShape(op, settings)
	ExpectEquals(t, err.Error(), "BlockArt: Shape svg string too long [M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8 h 37 M 8,0,0 q 8 L 4,4 l -4,4 v -8]")
}

func TestShapesOverlappingConcave(t *testing.T) {
	squareOut := Square1()
	squareIn := Square2()

	isOverlap := minerlib.IsShapesOverlapping(squareIn, squareOut) // Expect false
	ExpectFalse(t, isOverlap, "1) squareIn and squareOut. Expected false. Got true.\n")

	c := ConvexPolygon()

	isOverlap = minerlib.IsShapesOverlapping(squareOut, c) // Expect false
	ExpectFalse(t, isOverlap, "2) squareOut and c. Expected false. Got true.\n")

	isOverlap = minerlib.IsShapesOverlapping(squareIn, c) // Expect true
	ExpectTrue(t, isOverlap, "3) squareIn and c. Expected true. Got false.\n")

	c.Fill = minerlib.TRANSPARENT
	isOverlap = minerlib.IsShapesOverlapping(squareIn, c) // Expect false
	ExpectFalse(t, isOverlap, "4) squareIn and c transparent. Expected false. Got true.\n")
}

func TestShapesOverlappingConvex(t *testing.T) {
	square := Square3()
	triangleOut := Triangle1()
	triangleIn := Triangle2()

	isOverlap := minerlib.IsShapesOverlapping(square, triangleIn) // Expect true
	ExpectTrue(t, isOverlap, "1) square and triangleIn. Expected true. Got false.\n")

	triangleIn.Fill = minerlib.TRANSPARENT
	isOverlap = minerlib.IsShapesOverlapping(square, triangleIn) // Expect false
	ExpectFalse(t, isOverlap, "2) square and triangleIn. Expected false. Got true\n")

	isOverlap = minerlib.IsShapesOverlapping(square, triangleOut) // Expect false
	ExpectFalse(t, isOverlap, "3) square and triangleOut. Expected false. Got true\n")
}

func TestCirclesIntersecting(t *testing.T) {
	circle661 := Circle6_6_1()
	circle671 := Circle6_7_1()
	circle681 := Circle6_8_1()
	circle1061 := Circle10_6_1()
	convex := ConvexPolygon()

	isOverlap := minerlib.IsShapesOverlapping(circle661, circle1061) // Expect false
	ExpectFalse(t, isOverlap, "1) circle661 and circle1061. Expected false. Got true.\n")

	isOverlap = minerlib.IsShapesOverlapping(circle661, circle671) // Expect true
	ExpectTrue(t, isOverlap, "2) circle661 and circle671. Expected true. Got false.\n")

	isOverlap = minerlib.IsShapesOverlapping(circle661, convex)
	ExpectTrue(t, isOverlap, "3) circle661 and convex.\n")

	isOverlap = minerlib.IsShapesOverlapping(circle671, convex)
	ExpectTrue(t, isOverlap, "4) circle671 and convex.\n")

	isOverlap = minerlib.IsShapesOverlapping(circle661, circle681)
	ExpectTrue(t, isOverlap, "5) circle 661 and circle 681.\n")
}

func TestOtherCircleCases(t *testing.T) {
	square1 := Square1()
	circle1061 := Circle10_6_1()
	isOverlap := minerlib.IsShapesOverlapping(square1, circle1061)
	ExpectFalse(t, isOverlap, "1) square1 and circle 10,6 1.\n")

	// Circle in transparent polygon
	square4 := Square4()
	isOverlap = minerlib.IsShapesOverlapping(square4, circle1061)
	ExpectFalse(t, isOverlap, "2) square4 transparent and circle 10,6 1.\n")

	// Circle in filled polygon
	square4.Fill = "nontransparent"
	isOverlap = minerlib.IsShapesOverlapping(square4, circle1061)
	ExpectTrue(t, isOverlap, "3) square4 filled and circle 10,6 1.\n")

	// Polygon in transparent circle
	square2 := Square2()
	circle844 := Circle8_4_4()
	isOverlap = minerlib.IsShapesOverlapping(square2, circle844)
	ExpectFalse(t, isOverlap, "4) square2, transparent circle 8,4 4.\n")

	// Polygon in filled circle
	circle844.Fill = "nontransparent"
	isOverlap = minerlib.IsShapesOverlapping(square2, circle844)
	ExpectTrue(t, isOverlap, "5) square2, filled circle 8,4 4.\n")

	// Circle in transparent circle
	circle842 := Circle8_4_2()
	isOverlap = minerlib.IsShapesOverlapping(circle844, circle842)
	ExpectTrue(t, isOverlap, "6) circle 8,4 4 and circle 8,4 2.\n")

	// Circle in filled circle
	circle844.Fill = minerlib.TRANSPARENT
	isOverlap = minerlib.IsShapesOverlapping(circle844, circle842)
	ExpectFalse(t, isOverlap, "7) transparent circle 8,4 4 and circle 8,4 2.\n")
}

/*
type Operation struct {
	Type OperationType
	OperationNumber uint32
	Shape ShapeType
	Fill string // Can be "transparent" or "filled"
	Stroke string
	ShapeSVGString string
	ArtNodePubKey string
	ValidateBlockNum uint8
	ShapeHash string
	Nonce uint32
}
*/

func TestDrawAllShapes(t *testing.T) {
	convexPolygon := ConvexPolygon()
	convexPolygonOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "filled", "red", convexPolygon.ShapeToSVGPath(), "artnode0", 34, "convex_polygon", 129}
	squareOut := Square1()
	squareOutOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", squareOut.ShapeToSVGPath(), "artnode1", 34, "square_out", 129}
	squareIn := Square2()
	squareInOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", squareIn.ShapeToSVGPath(), "artnode2", 34, "square_in", 129}
	circleOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", "c 10,6 r 1", "artnode3", 34, "circle", 129}
	// TODO[sharon]: add and delete square5
	square5 := Square5()
	square5DrawOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", square5.ShapeToSVGPath(), "artnode2", 34, "square_in", 129}
	square5DeleteOp := Operation{blockartlib.DELETE, 2, blockartlib.PATH, "transparent", "red", square5.ShapeToSVGPath(), "artnode2", 34, "square_in", 129}

	operations := []Operation{convexPolygonOp, squareOutOp, squareInOp, circleOp, square5DrawOp, square5DeleteOp}
	settings := CanvasSettings{1024, 1024}
	validOps, invalidOps, _ := minerlib.DrawOperations(operations, settings)
	validString := ConcatOps(validOps)
	invalidString := ConcatOps(invalidOps)
	ExpectContains(t, validString, []string{"convex_polygon by artnode0", "square_out by artnode1"})
	ExpectContains(t, invalidString, []string{"square_in by artnode2", "circle by artnode3"})
}

func TestDrawAllShapesWithOwnership(t *testing.T) {
	convexPolygon := ConvexPolygon()
	convexPolygonOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "filled", "red", convexPolygon.ShapeToSVGPath(), "artnode0", 34, "convex_polygon", 129}
	squareOut := Square1()
	squareOutOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", squareOut.ShapeToSVGPath(), "artnode1", 34, "square_out", 129}
	squareIn := Square2()
	squareInOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", squareIn.ShapeToSVGPath(), "artnode0", 34, "square_in", 129}
	circleOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", "c 10,6 r 1", "artnode0", 34, "circle", 129}
	operations := []Operation{convexPolygonOp, squareOutOp, squareInOp, circleOp}
	settings := CanvasSettings{1024, 1024}
	validOps, invalidOps, _ := minerlib.DrawOperations(operations, settings)
	validString := ConcatOps(validOps)
	invalidString := ConcatOps(invalidOps)
	ExpectContains(t, validString, []string{"convex_polygon by artnode0", "square_out by artnode1", "square_in by artnode0", "circle by artnode0"})
	ExpectEquals(t, invalidString, "")
}

func Square1() Shape {
	so1 := LineSegment{Point{5, 2}, Point{6, 2}}
	so2 := LineSegment{Point{6, 2}, Point{6, 3}}
	so3 := LineSegment{Point{6, 3}, Point{5, 3}}
	so4 := LineSegment{Point{5, 3}, Point{5, 2}}
	soSides := []LineSegment{so1, so2, so3, so4}
	return Shape{"owner1", "square1", soSides, "transparent", "stroke", Point{0, 0}, 0}
}

func Square2() Shape {
	si1 := LineSegment{Point{8, 4}, Point{9, 4}}
	si2 := LineSegment{Point{9, 4}, Point{9, 5}}
	si3 := LineSegment{Point{9, 5}, Point{8, 5}}
	si4 := LineSegment{Point{8, 5}, Point{8, 4}}
	siSides := []LineSegment{si1, si2, si3, si4}
	return Shape{"owner2", "square2", siSides, "transparent", "stroke", Point{0, 0}, 0}
}

func Square3() Shape {
	s1 := LineSegment{Point{4, 3}, Point{5, 3}}
	s2 := LineSegment{Point{5, 3}, Point{5, 4}}
	s3 := LineSegment{Point{5, 4}, Point{4, 4}}
	s4 := LineSegment{Point{4, 4}, Point{4, 3}}
	sides := []LineSegment{s1, s2, s3, s4}
	return Shape{"owner3", "square3", sides, "transparent", "stroke", Point{0, 0}, 0}
}

func Square4() Shape {
	s1 := LineSegment{Point{15, 0}, Point{15, 10}}
	s2 := LineSegment{Point{15, 10}, Point{5, 10}}
	s3 := LineSegment{Point{5, 10}, Point{5, 0}}
	s4 := LineSegment{Point{5, 0}, Point{15, 0}}
	sides := []LineSegment{s1, s2, s3, s4}
	return Shape{"owner4", "square4", sides, "transparent", "stroke", Point{0, 0}, 0}
}

func Square5() Shape {
	s1 := LineSegment{Point{50, 50}, Point{50, 51}}
	s2 := LineSegment{Point{50, 51}, Point{49, 51}}
	s3 := LineSegment{Point{49, 51}, Point{49, 50}}
	s4 := LineSegment{Point{49, 50}, Point{50, 50}}
	sides := []LineSegment{s1, s2, s3, s4}
	return Shape{"owner13", "square5", sides, "transparent", "stroke", Point{0, 0}, 0}
}

func Triangle2() Shape {
	s1 := LineSegment{Point{6, 1}, Point{6, 6}}
	s2 := LineSegment{Point{6, 6}, Point{1, 4}}
	s3 := LineSegment{Point{1, 4}, Point{6, 1}}
	sides := []LineSegment{s1, s2, s3}
	return Shape{"owner4", "triangle1", sides, "filled", "stroke", Point{0, 0}, 0}
}

func Triangle1() Shape {
	s1 := LineSegment{Point{2, 3}, Point{3, 5}}
	s2 := LineSegment{Point{3, 5}, Point{1, 5}}
	s3 := LineSegment{Point{1, 5}, Point{2, 3}}
	sides := []LineSegment{s1, s2, s3}
	return Shape{"owner5", "triangle2", sides, "filled", "stroke", Point{0, 0}, 0}
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
	return Shape{"owner6", "concave_polygon", cSides, "solid", "stroke", Point{0, 0}, 0}
}

func Circle6_6_1() Shape {
	return Shape{"owner7", "circle 6,6 1", nil, "solid", "stroke", Point{6, 6}, 1}
}

func Circle6_7_1() Shape {
	return Shape{"owner8", "circle 6,7 1", nil, "solid", "stroke", Point{6, 7}, 1}
}

func Circle6_8_1() Shape {
	return Shape{"owner9", "circle 6,8 1", nil, "solid", "stroke", Point{6, 8}, 1}
}

func Circle10_6_1() Shape {
	return Shape{"owner10", "circle 10,6 1", nil, "solid", "stroke", Point{10, 6}, 1}
}

func Circle8_4_4() Shape {
	return Shape{"owner11", "circle 8,4 4", nil, "transparent", "stroke", Point{8, 4}, 4}
}

func Circle8_4_2() Shape {
	return Shape{"owner12", "circle 8,4 2", nil, "transparent", "stroke", Point{8, 4}, 2}
}