/*

Art app to be run to generate an HTML displaying the current
state of the blockchain in SVG strings.

Will wait X seconds and then run.

Usage:
go run html-art-app.go [minerAddress] [width] [height] [durationToWait]

Durations should look something like "10s" or "5ms" or "10m"
*/

package main

import (
	"blockartlib"
	"fmt"
	"os"
	"keys"
	"time"
)

var PATH_TO_FILE string = "canvas.html"

func main() {
	if(len(os.Args) < 4) {
		fmt.Println("Usage: go run html-art-app.go [MINER_ADDRESS] [width] [height] [duration]")
		os.Exit(2)
	}
	minerAddr := os.Args[1]
	width := os.Args[2]
	height:= os.Args[3]
	durationString := os.Args[4]
	durationToWait, err := time.ParseDuration(durationString)
	if err != nil {
		fmt.Println("Usage: go run html-art-app.go [MINER_ADDRESS] [width] [height] [duration]")
		fmt.Println("e.g. go run html-art-app.go 127.0.0.1:8000 1000 1000 10s")
		fmt.Println(err)
		os.Exit(2)
	}

	
	privateKey, _, err := keys.Generate()

	// Open a canvas.
	// TODO: use settings
	fmt.Printf("HTML-ART-APP: Calling OpenCanvas to Miner with address %s\n", minerAddr)
	canvas, _, err := blockartlib.OpenCanvas(minerAddr, *privateKey)
	if checkError(err) != nil {
		fmt.Println("HTML-ART-APP: there was an error opening the canvas", err)
		return
	}

    // Wait before requesting SVG strings so other miners can do work
	time.Sleep(durationToWait)

	// TODO: Get SVG strings from server and store appropriately
	/*
	fmt.Println("HTML-ART-APP: Getting Genesis block from Miner")
	blockHash, err := blockartlib.GetGenesisBlock()
	checkErr(err)
	fmt.Pritnln("HTML-ART-APP: Getting all SVG strings from the Miner")
	svgStrings, err := blockartlib.GetAllSvgStrings(blockhash)
	checkErr(err)
	*/
	var svgStrings []string = []string{
		"<path d=\"M5,2L6,2L6,3L5,3L5,2\" stroke=\"red\" fill=\"transparent\"/>",
	}
	// Generate HTML file from strings
    generateHTMLFile(width, height, svgStrings)

	// Close the canvas.
	fmt.Println("HTML-ART-APP: Closing the canvas")
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

func generateHTMLFile(width string, height string, svgStrings []string) {
	/* Given a bunch of SVG Strings, generate a viewable HTML file */
	if _, err := os.Stat(PATH_TO_FILE); err == nil {
		// If the file already exists, delete it
		err = os.Remove(PATH_TO_FILE)
		checkError(err)
	}

	// Open the file for writing
	file, err := os.OpenFile(PATH_TO_FILE, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	checkError(err)
	defer file.Close()


	// Write the HTML header, opening body, etc
	_, err = file.WriteString("<!DOCTYPE html>\n<html>\n<body>\n")
	checkError(err)
	_, err = file.WriteString(fmt.Sprintf("<svg width=\"%s\" height=\"%s\">\n", width, height))
    checkError(err)

	// Write SVG Lines
	for _, svgLine := range svgStrings {
	    file.WriteString(fmt.Sprintf("%s\n", svgLine))	
	}

	// Write closing tags
	_, err = file.WriteString("</svg>\n</body>\n</html>\n")
	checkError(err)

	fmt.Println("HTML-ART-APP: canvas.html generated")

}
