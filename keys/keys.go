package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"reflect"
)

func Generate() (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, err error) {
	privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Printf("Error while Generating Keys: %v\n", err)
		return nil, nil, err
	}
	publicKey = &privateKey.PublicKey
	prK := EncodePrivateKey(privateKey)
	fmt.Println(prK)
	return privateKey, publicKey, nil
}

func Encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	return EncodePrivateKey(privateKey), EncodePublicKey(publicKey)
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Error while Encoding Private Key: %v\n", err)
	}
	return hex.EncodeToString(x509Encoded)
}

func EncodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Printf("Error while Encoding Public Key: %v\n", err)
	}
	return hex.EncodeToString(x509EncodedPub)
}

func Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return DecodePrivateKey(pemEncoded), DecodePublicKey(pemEncodedPub)
}

func DecodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	hexDecodedBytes, err := hex.DecodeString(pemEncoded)
	if err != nil {
		fmt.Printf("Error while Dehexing Private Key: %v\n", err)
	}
	privateKey, err := x509.ParseECPrivateKey(hexDecodedBytes)
	if err != nil {
		fmt.Printf("Error while Decoding Private Key: %v\n", err)
	}
	return privateKey
}

func DecodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	hexDecodedBytes, err := hex.DecodeString(pemEncodedPub)
	if err != nil {
		fmt.Printf("Error while Dehexing Public Key: %v\n", err)
	}
	genericPublicKey, err := x509.ParsePKIXPublicKey(hexDecodedBytes)
	if err != nil {
		fmt.Printf("Error while Decoding Public Key: %v\n", err)
	}
	key := genericPublicKey.(*ecdsa.PublicKey)
	return key
}

func MatchPublicKeys(key0 *ecdsa.PublicKey, key1 *ecdsa.PublicKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}

func MatchPrivateKeys(key0 *ecdsa.PrivateKey, key1 *ecdsa.PrivateKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}

func MatchingPair(privateKey *ecdsa.PrivateKey, publickKey *ecdsa.PublicKey) (match bool) {
	return MatchPublicKeys(publickKey, &privateKey.PublicKey)
}
