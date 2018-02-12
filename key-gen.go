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
	fmt.Printf("%x\n", encPriv)
	fmt.Printf("%x\n", encPub)

	decodedPriv, decodedPub := keys.Decode(encPriv, encPub)

	if !keys.MatchPrivateKeys(privateKey, decodedPriv) {
		fmt.Println("Private keys do not match.")
	} else {
		fmt.Println("Do match")
	}
	if !keys.MatchPublicKeys(publicKey, decodedPub) {
		fmt.Println("Public keys do not match.")
	}
}
