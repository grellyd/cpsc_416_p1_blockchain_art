package minerlib

import (
	"blockartlib"
	"testing"
	"fmt"
	"keys"
)

var publicKeyStr = "3076301006072a8648ce3d020106052b81040022036200040426be3b44287dafed30cbb4b9bea7ecb9baf6910b4aaa70825fc604509a9bc36a2c6750638d1f44e6d95f1dfc68bc3d4e7a799f048d7019448f2b793a53c91129276a8c96d4ad7d58317cef1099b26f769671aa235071750e7b7511229b9f2c"

func TestMarshall(t *testing.T) {
	var tests = []struct {
		block    Block
		data     []byte
	}{
		{
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: keys.DecodePublicKey(publicKeyStr),
				Nonce: 4,
			},
			testBlock,
		},
	}
	for _, test := range tests {
		fmt.Printf("Testing Block: %v\n", test.block)
		data, err := test.block.MarshallBinary()
		if err != nil {
			t.Errorf("Bad Exit: \"TestMarshall(%v)\" produced err: %v", test, err)
		}
		if data == nil {
			t.Errorf("Bad Exit: No data, instead of %d", test.data)
		}
		err = nil
		for i, datum := range data {
			if datum != test.data[i] {
				err = fmt.Errorf("Error: Byte %d as '%d' doesn't match '%d'", i, datum, test.data[i])
				break
			}
		}
		if err != nil {
			t.Errorf("Bad Exit: Given \n'%d', instead of \n'%d':\n %v", data, test.data, err)
		}
	}
}

func TestMarshallErrors(t *testing.T) {
	// TODO: Test empty
}

func TestUnmarshall(t *testing.T) {
	var tests = []struct {
		data     []byte
		block    Block
	}{
		{
			testBlock,
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: keys.DecodePublicKey(publicKeyStr),
				Nonce: 4,
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
	var tests = []struct {
		block    Block
	}{
		{
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: keys.DecodePublicKey(publicKeyStr),
				Nonce: 4,
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
		if !keys.MatchPublicKeys(newBlock.MinerPublicKey, test.block.MinerPublicKey) {
			t.Errorf("Bad Exit: Keys Don't match! '%v' vs. '%v'", newBlock.MinerPublicKey, test.block.MinerPublicKey)
		}
		if newBlock.Nonce != test.block.Nonce {
			t.Errorf("Bad Exit: Nonces Don't match! '%d' vs. '%d'", newBlock.Nonce, test.block.Nonce)
		}
	}
}

var testBlock = []byte{80, 255, 129, 3, 1, 1, 5, 66, 108, 111, 99, 107, 1, 255, 130, 0, 1, 4, 1, 10, 80, 97, 114, 101, 110, 116, 72, 97, 115, 104, 1, 12, 0, 1, 10, 79, 112, 101, 114, 97, 116, 105, 111, 110, 115, 1, 255, 134, 0, 1, 14, 77, 105, 110, 101, 114, 80, 117, 98, 108, 105, 99, 75, 101, 121, 1, 255, 136, 0, 1, 5, 78, 111, 110, 99, 101, 1, 6, 0, 0, 0, 39, 255, 133, 2, 1, 1, 24, 91, 93, 42, 98, 108, 111, 99, 107, 97, 114, 116, 108, 105, 98, 46, 79, 112, 101, 114, 97, 116, 105, 111, 110, 1, 255, 134, 0, 1, 255, 132, 0, 0, 255, 135, 255, 131, 3, 1, 2, 255, 132, 0, 1, 9, 1, 4, 84, 121, 112, 101, 1, 4, 0, 1, 15, 79, 112, 101, 114, 97, 116, 105, 111, 110, 78, 117, 109, 98, 101, 114, 1, 4, 0, 1, 12, 79, 112, 101, 114, 97, 116, 105, 111, 110, 83, 105, 103, 1, 12, 0, 1, 5, 83, 104, 97, 112, 101, 1, 4, 0, 1, 4, 70, 105, 108, 108, 1, 12, 0, 1, 6, 83, 116, 114, 111, 107, 101, 1, 12, 0, 1, 14, 83, 104, 97, 112, 101, 83, 86, 71, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 13, 65, 114, 116, 78, 111, 100, 101, 80, 117, 98, 75, 101, 121, 1, 12, 0, 1, 5, 78, 111, 110, 99, 101, 1, 6, 0, 0, 0, 47, 255, 135, 3, 1, 1, 9, 80, 117, 98, 108, 105, 99, 75, 101, 121, 1, 255, 136, 0, 1, 3, 1, 5, 67, 117, 114, 118, 101, 1, 16, 0, 1, 1, 88, 1, 255, 138, 0, 1, 1, 89, 1, 255, 138, 0, 0, 0, 10, 255, 137, 5, 1, 2, 255, 140, 0, 0, 0, 121, 255, 130, 1, 10, 112, 97, 114, 101, 110, 116, 72, 97, 115, 104, 2, 1, 21, 42, 101, 108, 108, 105, 112, 116, 105, 99, 46, 67, 117, 114, 118, 101, 80, 97, 114, 97, 109, 115, 255, 141, 3, 1, 1, 11, 67, 117, 114, 118, 101, 80, 97, 114, 97, 109, 115, 1, 255, 142, 0, 1, 7, 1, 1, 80, 1, 255, 138, 0, 1, 1, 78, 1, 255, 138, 0, 1, 1, 66, 1, 255, 138, 0, 1, 2, 71, 120, 1, 255, 138, 0, 1, 2, 71, 121, 1, 255, 138, 0, 1, 7, 66, 105, 116, 83, 105, 122, 101, 1, 4, 0, 1, 4, 78, 97, 109, 101, 1, 12, 0, 0, 0, 254, 1, 122, 255, 142, 254, 1, 11, 1, 49, 2, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 1, 49, 2, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 199, 99, 77, 129, 244, 55, 45, 223, 88, 26, 13, 178, 72, 176, 167, 122, 236, 236, 25, 106, 204, 197, 41, 115, 1, 49, 2, 179, 49, 47, 167, 226, 62, 231, 228, 152, 142, 5, 107, 227, 248, 45, 25, 24, 29, 156, 110, 254, 129, 65, 18, 3, 20, 8, 143, 80, 19, 135, 90, 198, 86, 57, 141, 138, 46, 209, 157, 42, 133, 200, 237, 211, 236, 42, 239, 1, 49, 2, 170, 135, 202, 34, 190, 139, 5, 55, 142, 177, 199, 30, 243, 32, 173, 116, 110, 29, 59, 98, 139, 167, 155, 152, 89, 247, 65, 224, 130, 84, 42, 56, 85, 2, 242, 93, 191, 85, 41, 108, 58, 84, 94, 56, 114, 118, 10, 183, 1, 49, 2, 54, 23, 222, 74, 150, 38, 44, 111, 93, 158, 152, 191, 146, 146, 220, 41, 248, 244, 29, 189, 40, 154, 20, 124, 233, 218, 49, 19, 181, 240, 184, 192, 10, 96, 177, 206, 29, 126, 129, 157, 122, 67, 29, 124, 144, 234, 14, 95, 1, 254, 3, 0, 1, 5, 80, 45, 51, 56, 52, 0, 1, 49, 2, 4, 38, 190, 59, 68, 40, 125, 175, 237, 48, 203, 180, 185, 190, 167, 236, 185, 186, 246, 145, 11, 74, 170, 112, 130, 95, 198, 4, 80, 154, 155, 195, 106, 44, 103, 80, 99, 141, 31, 68, 230, 217, 95, 29, 252, 104, 188, 61, 1, 49, 2, 78, 122, 121, 159, 4, 141, 112, 25, 68, 143, 43, 121, 58, 83, 201, 17, 41, 39, 106, 140, 150, 212, 173, 125, 88, 49, 124, 239, 16, 153, 178, 111, 118, 150, 113, 170, 35, 80, 113, 117, 14, 123, 117, 17, 34, 155, 159, 44, 0, 1, 4, 0}
