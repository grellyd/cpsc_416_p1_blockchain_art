package minerlib

import (
	"blockartlib"
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func GenerateKeys() (publicKey *ecdsa.PublicKey, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to Generate Keys: %v", err)
	}
	return privateKey.Public().(*ecdsa.PublicKey), nil
}


func TestMarshall(t *testing.T) {
	publicKey, err := GenerateKeys()
	if err != nil {
		t.Errorf("Bad Exit: \"TestMarshall()\" produced err: %v", err)
	}
	var tests = []struct {
		block    Block
		data     []byte
	}{
		// TODO: Add more testing scenarios
		{
			Block{},
			[]byte{},
		},
		{
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: publicKey,
				nonce: 4,
			},
			[]byte{},
		},
	}
	for _, test := range tests {
		data, err := test.block.MarshallBinary()
		if err != nil {
			t.Errorf("Bad Exit: \"TestMarshall(%v)\" produced err: %v", test, err)
		}
		if data == nil {
			t.Errorf("Bad Exit: No data, instead of %d", test.data)
		}
		// TODO: Figure out what the data should actually look like
	}
}

func TestUnmarshall(t *testing.T) {
	publicKey, err := GenerateKeys()
	if err != nil {
		t.Errorf("Bad Exit: \"TestUnmarshall()\" produced err: %v", err)
	}
	var tests = []struct {
		data     []byte
		block    Block
	}{
		// TODO: Get some test byte arrays
		{
			[]byte{},
			Block{},
		},
		{
			[]byte{},
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: publicKey,
				nonce: 4,
			},
		},
	}
	for _, test := range tests {
		block, err := UnmarshallBinary(test.data)
		if err != nil {
			t.Errorf("Bad Exit: \"TestUnmarshall(%v)\" produced err: %v", test, err)
		}
		if block == nil {
			t.Errorf("Bad Exit: No block produced, instead of %d", test.block)
		}
	}
}

func TestMarshallUnMarshall(t *testing.T) {
	publicKey, err := GenerateKeys()
	if err != nil {
		t.Errorf("Bad Exit: \"TestMarshallUnmarshall()\" produced err: %v", err)
	}
	var tests = []struct {
		block    Block
	}{
		{
			Block{},
		},
		{
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: publicKey,
				nonce: 4,
			},
		},
	}
	for _, test := range tests {
		data, err := test.block.MarshallBinary()
		if err != nil {
			t.Errorf("Bad Exit: \"TestMarshallUnmarshall(%v)\" produced err: %v", test, err)
		}
		newBlock, err := UnmarshallBinary(data)
		if err != nil {
			t.Errorf("Bad Exit: \"TestMarshallUnmarshall(%v)\" produced err: %v", test, err)
		}
		if newBlock == nil || data == nil { 
			t.Errorf("Bad Exit: No results or data!")
		}
		if newBlock.ParentHash != test.block.ParentHash {
			t.Errorf("Bad Exit: Parents Don't match! '%s' vs. '%s'", newBlock.ParentHash, test.block.ParentHash)
		}
		for i, op := range newBlock.Operations {
			if op != test.block.Operations[i] {
				t.Errorf("Bad Exit: Operation addresses Don't match! '%v' vs. '%v'", op, test.block.Operations[i])
			}
		}
		if newBlock.MinerPublicKey != test.block.MinerPublicKey {
			t.Errorf("Bad Exit: Keys Don't match! '%v' vs. '%v'", newBlock.MinerPublicKey, test.block.MinerPublicKey)
		}
		if newBlock.nonce != test.block.nonce {
			t.Errorf("Bad Exit: Nonces Don't match! '%d' vs. '%d'", newBlock.nonce, test.block.nonce)
		}
	}
}
