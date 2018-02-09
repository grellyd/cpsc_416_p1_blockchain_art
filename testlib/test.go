package main

import (
	"fmt"
	"blockartlib"
	"minerlib"
)

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

func main() {
	validSvg := "M 0,1 L 4,11.2 V 8 H 9.72123"
	validOp := blockartlib.Operation{1, 0, "sig", blockartlib.ShapeType(1), "transparent", "red", validSvg, "pubkey", 9}
	invalidSvg := "89z"
	invalidOp := blockartlib.Operation{1, 0, "sig", blockartlib.ShapeType(1), "transparent", "red", invalidSvg, "pubkey", 9}
	validShape, _ := minerlib.OperationToShape(validOp)
	fmt.Printf("valid shape: %v\n", validShape)
	invalidShape, _ := minerlib.OperationToShape(invalidOp)
	fmt.Printf("invalid shape: %v\n", invalidShape)
}