package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strings"

	"./minerlib"
)

func main() {
	GetKeyPair() // Temporary, keys will be passed in as command line args
	// Need to print then pass to client
	// Connect to server
	localIP := "127.0.0.1:0"
	//serverAddr := "tbd"
	var nbrs [256]int
	m := minerlib.Miner{"server addr", nbrs, 0,localIP, false}

	fmt.Printf("miner ip: %v, m: %v, \n", localIP, m)
}



func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func GetKeyPair() {
	curve := elliptic.P256()
	r := strings.NewReader("Hello, Reader!")
	keys, _ := ecdsa.GenerateKey(curve, r)
	fmt.Printf("Keys: %v\n", keys.PublicKey)
}
