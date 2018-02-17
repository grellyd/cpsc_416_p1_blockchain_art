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
	"keys"
	"os"
)

func main() {
	validateNum := 0 // TODO: Change this to a bigger number for submission
	shapes := []string{}
	blocks := []string{}
	minerAddr := os.Args[1]
	privateKey, _, err := keys.Generate()
	fmt.Printf("%v\n", privateKey.PublicKey)
	colour := "black"

	// Open a canvas.
	// TODO: use settings
	fmt.Printf("ART-APP: Calling OpenCanvas to Miner with address %s\n", minerAddr)
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privateKey)
	fmt.Println("ART-APP: Canvas is ", canvas)
	if checkError(err) != nil {
		fmt.Println("ART-APP: there was an error opening the canvas", err)
		return
	}

	// Draw squares along bottom
	for i := 0; i < 5000; i++ {
		fmt.Printf("ART-APP: Drawing square at x = %v.\n", i)
		svg := fmt.Sprintf("M %v , 1 h 1 v 1 h -1 v -1", i, i)
		shapeHash, blockHash, _, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, svg, "filled", colour)
		if checkError(err) != nil {
			fmt.Printf("ART-APP: There was an error with calling AddShape: \n")
			fmt.Println(err)
		}
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}

	fmt.Println("Closing the canvas")
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
