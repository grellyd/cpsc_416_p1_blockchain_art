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
	colour := "red"

	fmt.Printf("[art-app] shapes; %v\n", shapes)
	// Open a canvas.
	// TODO: use settings
	fmt.Printf("ART-APP: Calling OpenCanvas to Miner with address %s\n", minerAddr)
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privateKey)
	fmt.Println("ART-APP: Canvas is ", canvas)
	if checkError(err) != nil {
		fmt.Println("ART-APP: there was an error opening the canvas", err)
		return
	}

	// Add a line.
	fmt.Println("ART-APP: Calling AddShape to add a red line")
	shapeHash, blockHash, ink, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M 0,0 L 0, 5", "transparent", "red")
	fmt.Printf("ART-APP: shapeHash: '%v'\n", shapeHash)
	if checkError(err) != nil {
		fmt.Printf("ART-APP: There was an error with calling AddShape: \n")
		fmt.Println(err)
	}
	shapes = append(shapes, shapeHash)
	fmt.Printf("[art-app] shapes; %v\n", shapes)
	blocks = append(blocks, blockHash)

	// Add another line.
	shapeHash, blockHash, ink2, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M 0,0 L 5 ,0", "transparent", "blue")
	fmt.Printf("ART-APP: shapeHash: '%v'\n", shapeHash)
	fmt.Println("ART-APP: Calling AddShape to add a filled circle. Intersects with polygon. Gets drawn.")
	shapeHash, blockHash, ink2, err = canvas.AddShape(uint8(validateNum), blockartlib.CIRCLE, "c 10,6 r 1", "filled", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	}
	if ink2 <= ink {
		checkError(fmt.Errorf("Err! ink2 not > ink1"))
	}
	shapes = append(shapes, shapeHash)
	fmt.Printf("[art-app] shapes: %v\n", shapes)
	fmt.Printf("[art-app] shapes[0]: %v\n", shapes[0])
	blocks = append(blocks, blockHash)

	// Delete the first line.
	fmt.Println("ART-APP: Deleting the first line")
	ink3, err := canvas.DeleteShape(uint8(validateNum), shapes[0])
	if checkError(err) != nil {
		fmt.Println(err)
	}

	// assert ink3 > ink2
	if ink3 <= ink2 {
		checkError(fmt.Errorf("err! ink3 not > ink4"))
	}

	// Draw square in transparent circle.
	fmt.Println("ART-APP: Will draw transparent circle then filled square inside.")
	shapeHash, blockHash, _, err = canvas.AddShape(uint8(validateNum), blockartlib.CIRCLE, "c 50, 50 r 10", "transparent", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	} else {
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}

	shapeHash, blockHash, ink5, err := canvas.AddShape(uint8(validateNum), blockartlib.CIRCLE, "M50,50 h 1 v -1 h -1 v 1", "transparent", colour)
	if checkError(err) != nil {
		fmt.Println(err)
	} else {
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}

	fmt.Println("ART-APP: Drawing line that intersects with circle 50, 50 r 10.")
	shapeHash, blockHash, ink6, err := canvas.AddShape(uint8(validateNum), blockartlib.PATH, "M50,50 h 60", "transparent", colour)
	if _, ok := err.(*blockartlib.ShapeOverlapError); ok {
		fmt.Printf("Got ShapeOverlapError as expected. Err: %v\n", err)
	} else {
		shapes = append(shapes, shapeHash)
		blocks = append(blocks, blockHash)
	}

	if ink6 <= ink5 {
		checkError(fmt.Errorf("err! ink5 not > ink6"))
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
