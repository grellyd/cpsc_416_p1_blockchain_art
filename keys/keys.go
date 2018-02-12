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
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)

	var b []byte = x509Encoded
	var c []byte = x509EncodedPub
	return string(b), string(c)

}

func Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, _ := x509.ParseECPrivateKey([]byte(pemEncoded))
	genericPublicKey, _ := x509.ParsePKIXPublicKey([]byte(pemEncodedPub))
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

func MatchPublicKeys(key0 *ecdsa.PublicKey, key1 *ecdsa.PublicKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}

func MatchPrivateKeys(key0 *ecdsa.PrivateKey, key1 *ecdsa.PrivateKey) (match bool) {
	return reflect.DeepEqual(key0, key1)
}
