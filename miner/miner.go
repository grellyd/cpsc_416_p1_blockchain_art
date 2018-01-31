package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strings"

	"../minerlib"
)

func main() {
	curve := elliptic.P256()
	r := strings.NewReader("Hello, Reader!")

	keys, err := ecdsa.GenerateKey(curve, r)
	CheckError(err)
	// Connect to server
	localIP := "127.0.0.1:0"
	//serverAddr := "tbd"
	var nbrs [256]int
	m := minerlib.Miner{"server addr", nbrs, localIP}

	fmt.Printf("miner ip: %v, m: %v, keys.pub %v\n", localIP, m, keys.PublicKey)
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
