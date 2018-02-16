/*

A trivial application to illustrate how the blockartlib library can be
used from an application in project 1 for UBC CS 416 2017W2.

Usage:
go run art-app.go
*/

package main

import (
	"blockartlib"
	"fmt"
	"os"
	"keys"
)

func main() {
	shapes := []string{}
	blocks := []string{}
	minerAddr := os.Args[1]
	privateKey, _, err := keys.Generate()
	fmt.Printf("%v\n", privateKey.PublicKey)

	// Open a canvas.
	// TODO: use settings
	 canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privateKey)
	fmt.Println("Canvas", canvas, "error ", err)
	if checkError(err) != nil {
		return
	}

	validateNum := 2

	// Add a line.
	shapeHash, blockHash, ink, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M 0 0 L 0 5", "transparent", "red")
	if checkError(err) != nil {
		return
	}
	shapes = append(shapes, shapeHash)
	blocks = append(blocks, blockHash)

	// Add another line.
	shapeHash, blockHash, ink2, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M 0 0 L 5 0", "transparent", "blue")
	if checkError(err) != nil {
		return
	}
	if ink2 <= ink {
		checkError(fmt.Errorf("Err! ink2 not > ink1"))
	}
	shapes = append(shapes, shapeHash)
	blocks = append(blocks, blockHash)

	// Delete the first line.
	ink3, err := canvas.DeleteShape(uint8(validateNum), shapeHash)
	if checkError(err) != nil {
		return
	}

	// assert ink3 > ink2
	if ink3 <= ink2 {
		checkError(fmt.Errorf("Err! ink3 not > ink4"))
	}

	// Close the canvas.
	_, err = canvas.CloseCanvas()
	if checkError(err) != nil {
		return
	}
}

// If error is non-nil, print it out and return it.
func checkError(err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		return err
	}
	return nil
}
