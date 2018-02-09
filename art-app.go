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
	//"crypto/ecdsa"
	//"strconv"
	"fmt"
	"os"
	"encoding/hex"
	"crypto/x509"
)

func main() {
	shapes := []string{}
	blocks := []string{}
	//minerAddr := "127.0.0.1:8080"
	minerAddr := "127.0.0.1:65367"
	// TODO: use crypto/ecdsa to read pub/priv keys from a file argument.
	pk := "3081a40201010430db825fb603900f7e9de3e9323df1be08b9d2fd16418a3711638e97004570c99a4cc5e8edba8be71f7f8aa0e295b44a0aa00706052b81040022a164036200043b70a5f19bd65e8f401042363261379bba5c7cc4772ba15a30fb6eee3a8731166478f735c1cd1583237d7052701fe83e2c7feda965ab2cc02800e1a2a463d683fcd68d7280650086b5b31167435581c5c31cdd420328b012669406fe0928b7fe"
	//puk := "3076301006072a8648ce3d020106052b81040022036200043b70a5f19bd65e8f401042363261379bba5c7cc4772ba15a30fb6eee3a8731166478f735c1cd1583237d7052701fe83e2c7feda965ab2cc02800e1a2a463d683fcd68d7280650086b5b31167435581c5c31cdd420328b012669406fe0928b7fe"
	c, _ := hex.DecodeString(pk)
	//d, _ := hex.DecodeString(puk)
	privKey, _ := x509.ParseECPrivateKey(c)
	//genericPublicKey, _ := x509.ParsePKIXPublicKey(d)
	//publicKey := genericPublicKey.(*ecdsa.PublicKey)

	// Open a canvas.
	// TODO: use settings
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privKey)
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
