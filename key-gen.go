package main

import (
	"fmt"
	"keys"
)

func main()  {
	privateKey, publicKey, err := keys.Generate()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	encPriv, encPub := keys.Encode(privateKey, publicKey)
	fmt.Printf("Private:\n------\n%x\n\n", encPriv)
	fmt.Printf("Public:\n------\n%x\n\n", encPub)
	fmt.Println("-----")

	decodedPriv, decodedPub := keys.Decode(encPriv, encPub)

	if !keys.MatchPrivateKeys(privateKey, decodedPriv) {
		fmt.Println("Private keys do not match.")
	} else {
		fmt.Println("Keys match")
	}
	if !keys.MatchPublicKeys(publicKey, decodedPub) {
		fmt.Println("Public keys do not match.")
	}
}
