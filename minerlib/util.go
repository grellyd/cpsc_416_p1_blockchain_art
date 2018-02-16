package minerlib

import (
	"blockartlib"
	"strconv"
)

func OpToSvg(op Operation, settings CanvasSettings) (svg string) {
	svg = "<svg>"
	if op.Shape == blockartlib.PATH {
		svg += "<path d=\"" + op.ShapeSVGString + "\" "
		svg += "stroke=\"" + op.Stroke + "\" "
		svg += "fill =\"" + op.Fill + "\"/>"
	} else if op.Shape == blockartlib.CIRCLE {
		c, _ := OperationToShape(op, settings)
		svg += "<circle cx=\"" + strconv.FormatFloat(c.Center.X, 'f', -1, 64) + "\" "
		svg += "cy = \"" + strconv.FormatFloat(c.Center.Y, 'f', -1, 64) + "\" "
		svg += "r = \"" + strconv.FormatFloat(c.Radius, 'f', -1, 64) + "\" "
		svg += "stroke=\"" + op.Stroke + "\" "
		svg += "fill =\"" + op.Fill + "\"/>"
	} else {
		svg = ""
	}
	svg += "</svg>"
	return svg
}


