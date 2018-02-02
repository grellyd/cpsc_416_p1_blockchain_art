package blockartlib

// Canvas data structure that validates additions/removals
// from the canvas

type CanvasData struct {
	Points [][]int // TODO: Need to update the dimensions when creating canvas
	Operations map[string]Operation
}

func (cd *CanvasData) ValidateAddShape() (b bool) {
	return b
}