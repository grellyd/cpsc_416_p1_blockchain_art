/*
Mining is a package to mine a valid hash given a byte array and difficulty (number of zeros)
*/

package mine

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"math/rand"
	"fmt"
	"encoding/binary"
	"strings"
)

var usedNonces map[uint32]bool

func Data(data []byte, difficulty uint8) (nonce uint32, err error) {
	usedNonces = make(map[uint32]bool)
	for {
		n, err := NewNonce()
		if err != nil {
			return 0, fmt.Errorf("Unable to mine nonce for data: %v", err)
		}
		if Valid(MD5Hash(data, n), difficulty) {
			return n, nil
		}
	}
}

func Valid(hash string, difficulty uint8) (valid bool) {
	last_index := strings.LastIndex(hash, zeroString(difficulty))
	if last_index == len(hash)-int(difficulty) {
		valid = true
	}
	return valid
}

func MD5Hash(data []byte, nonce uint32) (hash string) {
	h := md5.New()
	h.Write(append(data, asByteArr(nonce)...))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

func asByteArr(val uint32) []byte {
	a := make([]byte, 4)
    binary.LittleEndian.PutUint32(a, val)
	return a
}

func NewNonce() (newNonce uint32, err error) {
	// TODO check for all tries
	for {
		newNonce = randomNonce()
		if !usedNonces[newNonce] {
			usedNonces[newNonce] = true
			return newNonce, nil
		}
	}
}

func randomNonce() (nonce uint32) {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func zeroString(num uint8) string {
	str := ""
	for i := uint8(0); i < num; i++ {
		str += "0"
	}
	return str
}
