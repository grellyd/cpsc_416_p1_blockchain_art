/*

A trivial application to illustrate how the blockartlib library can be
used from an application in project 1 for UBC CS 416 2017W2.

Usage:
go run art-app.go
*/

package main

import (
	"time"
	"blockartlib"
	"fmt"
	"os"
	"keys"
)

func main() {
	time.Sleep(30,000,000,000) // Sleeps for 30 seconds. Unit is nanoseconds
	validateNum := 0 // TODO: Change this to a bigger number for submission
	shapes := []string{}
	blocks := []string{}
	minerAddr := os.Args[1]
	privateKey, _, err := keys.Generate()
	fmt.Printf("%v\n", privateKey.PublicKey)
	colour := "blue"

	// Open a canvas.
	// TODO: use settings
	fmt.Printf("ART-APP1: Calling OpenCanvas to Miner with address %s\n", minerAddr)
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privateKey)
	fmt.Println("ART-APP1: Canvas is ", canvas)
	if checkError(err) != nil {
		fmt.Println("ART-APP1: there was an error opening the canvas", err)
		return
	}

	// Add a line.
	fmt.Println("ART-APP1: Calling AddShape to add a red filled square")
	shapeHash, blockHash, ink, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M8,4L9,4L9,5L8,5L8,4", "filled", colour)
	if checkError(err) != nil {
		fmt.Printf("ART-APP1: There was an error with calling AddShape: \n")
		fmt.Println(err)
		return
	}
	shapes = append(shapes, shapeHash)
	blocks = append(blocks, blockHash)

	// Add another line.
    fmt.Println("ART-APP1: Calling AddShape to add a transparent triangle. Expect fail")
	shapeHash, blockHash, ink2, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M6,1L6,6L1,4L6,1", "filled", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	} else {
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}
	if ink2 <= ink {
		checkError(fmt.Errorf("Err! ink2 not > ink1"))
	}

	// Delete the first line.
	fmt.Println("ART-APP1: Deleting the first line")
	ink3, err := canvas.DeleteShape(uint8(validateNum), shapeHash)
	if checkError(err) != nil {
		fmt.Println(err)
	} 

	// assert ink3 > ink2
	if ink3 <= ink2 {
		checkError(fmt.Errorf("Err! ink3 not > ink4"))
	}

	fmt.Println("ART-APP1: Drawing square that intersects with ART-APP's polygon")
	shapeHash, blockHash, ink4, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M4,3 h 1 v 1 h -1 v -1", "filled", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	} else {
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}

	fmt.Println("ART-APP1: Drawing shape with out of bounds svg string") // M -4,-3 invalid
	shapeHash, blockHash, ink4, err = canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M-4,-3 h 1 v 1 h -1 v -1", "filled", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	} 

	fmt.Println("ART-APP1: Drawing shape with invalid svg string") // Q not supported
	shapeHash, blockHash, ink4, err = canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M-4,-3 Q 1 v 1 h -1 v -1", "filled", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	}

	fmt.Println("Closing the canvas")
	// Close the canvas.
	_, err = canvas.CloseCanvas()
	if checkError(err) != nil {
		return
	}
}
