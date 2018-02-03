/*

A trivial application to illustrate how the blockartlib library can be
used from an application in project 1 for UBC CS 416 2017W2.

Usage:
go run art-app.go
*/

package main

// Expects blockartlib.go to be in the ./blockartlib/ dir, relative to
// this art-app.go file
import "blockartlib"

import (
	"crypto/ecdsa"
	//"strconv"
	"fmt"
	"os"
)

func main() {
	shapes := []string{}
	blocks := []string{}
	minerAddr := "127.0.0.1:8080"
	// TODO: use crypto/ecdsa to read pub/priv keys from a file argument.
	privKey := ecdsa.PrivateKey{}
/*	args := os.Args[1:]
	var privKey ecdsa.PrivateKey{} = args[2]
	privKey.PublicKey = mar*/
	// Open a canvas.
	// TODO: use settings
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, privKey)
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
