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
	"keys"
	"os"
	"time"
	"minerlib"
)

var PATH_TO_FILE string = "canvas.html"

type TreeNode struct {
	Hash          string
	Children      []*TreeNode
	Depth         int
	ParentsInPath []string
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run html-art-app.go [MINER_ADDRESS] [width] [height] [duration]")
		os.Exit(2)
	}

	minerAddr := os.Args[1]
	width := os.Args[2]
	height := os.Args[3]
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

	var svgStrings []string = []string{
		"<path d=\"M5,2L6,2L6,3L5,3L5,2\" stroke=\"red\" fill=\"transparent\"/>",
	}

	genNode := TreeNode{blockartlib.GetGenesisBlock(), []*TreeNode{}, 0, []string()}
	GetTree(&genNode)
	deepest := DeepestNode(&genNode)
	for _, i := range deepest.ParentsInPath {
		shapehashes, err := blockartlib.GetShapes(i)
		for _, h := range shapehashes {
			svgStrings, _ = append(svgStrings, blockartlib.GetSvgString(h))
		}
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

/*
	GetGenesisBlock() (blockHash string, err error)

	// Retrieves the children blocks of the block identified by blockHash.
	// Can return the following errors:
	// - DisconnectedError
	// - InvalidBlockHashError
	GetChildren(blockHash string) (blockHashes []string, err error)
*/
func GetTree(genNode *TreeNode) {
	genNode.Children, err = blockartlib.GetChildren(genNode.Hash)
	if err != nil {
		return
	}
	if len(genNode.Children) <= 0 {
		return
	}
	for _, child := range genNode.Children {
		GetTree(*TreeNode{child, []string{}, genNode.Depth + 1, append(genNode.ParentsInPath, genNode.Hash)})
	}
}

func DeepestNode(genNode *TreeNode) (res *TreeNode) {
	Find(genNode, res, 0, -1)
	return res
}

func Find(genNode, res *TreeNode, level, maxLevel *int) {
	if genNode != nil {
		for _, child := range genNode.Children {
			level++
			Find(child, res, level, *maxLevel)
			if level > maxLevel {
				res = &genNode
			}
		}
	}
}
