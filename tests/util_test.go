package tests

import (
	"blockartlib"
	"minerlib"
	"testing"
)

func TestOpToString(t *testing.T) {
	circleOp := Operation{blockartlib.DRAW, 2, blockartlib.CIRCLE, "transparent", "red", "c 10,6 r 1", "artnode0", 34, "circle", 129}
	settings := blockartlib.CanvasSettings{1000, 1000}
	s := minerlib.OpToSvg(circleOp, settings)
	ExpectEquals(t, s, "<circle cx=\"10\" cy = \"6\" r = \"1\" stroke=\"red\" fill =\"transparent\"/>")

	squareOut := Square1()
	squareOutOp := Operation{blockartlib.DRAW, 2, blockartlib.PATH, "transparent", "red", squareOut.ShapeToSVGPath(), "artnode1", 34, "square_out", 129}
	s = minerlib.OpToSvg(squareOutOp, settings)
	ExpectEquals(t, s, "<path d=\"M5,2L6,2L6,3L5,3L5,2\" stroke=\"red\" fill =\"transparent\"/>")
}
