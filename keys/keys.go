package keys

import (
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/elliptic"
	"crypto/rand"
	"reflect"
)

func Generate() (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, err error) {
	privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey = &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func Encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	return EncodePrivateKey(privateKey), EncodePublicKey(publicKey)
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) (string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	var b []byte = x509Encoded
	return string(b)
}

func EncodePublicKey(publicKey *ecdsa.PublicKey) (string) {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)

	var c []byte = x509EncodedPub
	return string(c)
}

func Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return DecodePrivateKey(pemEncoded), DecodePublicKey(pemEncodedPub)
}

func DecodePrivateKey(pemEncoded string) (*ecdsa.PrivateKey) {
	privateKey, _ := x509.ParseECPrivateKey([]byte(pemEncoded))
	return privateKey
}

func DecodePublicKey(pemEncodedPub string) (*ecdsa.PublicKey) {
	genericPublicKey, _ := x509.ParsePKIXPublicKey([]byte(pemEncodedPub))
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}


func MatchPublicKeys(key0 *ecdsa.PublicKey, key1 *ecdsa.PublicKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}

func MatchPrivateKeys(key0 *ecdsa.PrivateKey, key1 *ecdsa.PrivateKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}
