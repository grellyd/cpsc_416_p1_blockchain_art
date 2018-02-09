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
			Block{
				ParentHash: "parentHash",
				Operations: []*blockartlib.Operation{},
				MinerPublicKey: publicKey,
				Nonce: 4,
			},
			[]byte{80 ,255 ,129 ,3 ,1 ,1 ,5 ,66 ,108 ,111 ,99 ,107 ,1 ,255 ,130 ,0 ,1 ,4 ,1 ,10 ,80 ,97 ,114 ,101 ,110 ,116 ,72 ,97 ,115 ,104 ,1 ,12 ,0 ,1 ,10 ,79 ,112 ,101 ,114 ,97 ,116 ,105 ,111 ,110 ,115 ,1 ,255 ,134 ,0 ,1 ,14 ,77 ,105 ,110 ,101 ,114 ,80 ,117 ,98 ,108 ,105 ,99 ,75 ,101 ,121 ,1 ,255 ,136 ,0 ,1 ,5 ,78 ,111 ,110 ,99 ,101 ,1 ,6 ,0 ,0 ,0 ,39 ,255 ,133 ,2 ,1 ,1 ,24 ,91 ,93 ,42 ,98 ,108 ,111 ,99 ,107 ,97 ,114 ,116 ,108 ,105 ,98 ,46 ,79 ,112 ,101 ,114 ,97 ,116 ,105 ,111 ,110 ,1 ,255 ,134 ,0 ,1 ,255 ,132 ,0 ,0 ,255 ,135 ,255 ,131 ,3 ,1 ,2 ,255 ,132 ,0 ,1 ,9 ,1 ,4 ,84 ,121 ,112 ,101 ,1 ,4 ,0 ,1 ,15 ,79 ,112 ,101 ,114 ,97 ,116 ,105 ,111 ,110 ,78 ,117 ,109 ,98 ,101 ,114 ,1 ,4 ,0 ,1 ,12 ,79 ,112 ,101 ,114 ,97 ,116 ,105 ,111 ,110 ,83 ,105 ,103 ,1 ,12 ,0 ,1 ,5 ,83 ,104 ,97 ,112 ,101 ,1 ,4 ,0 ,1 ,4 ,70 ,105 ,108 ,108 ,1 ,12 ,0 ,1 ,6 ,83 ,116 ,114 ,111 ,107 ,101 ,1 ,12 ,0 ,1 ,14 ,83 ,104 ,97 ,112 ,101 ,83 ,86 ,71 ,83 ,116 ,114 ,105 ,110 ,103 ,1 ,12 ,0 ,1 ,13 ,65 ,114 ,116 ,78 ,111 ,100 ,101 ,80 ,117 ,98 ,75 ,101 ,121 ,1 ,12 ,0 ,1 ,5 ,78 ,111 ,110 ,99 ,101 ,1 ,6 ,0 ,0 ,0 ,47 ,255 ,135 ,3 ,1 ,1 ,9 ,80 ,117 ,98 ,108 ,105 ,99 ,75 ,101 ,121 ,1 ,255 ,136 ,0 ,1 ,3 ,1 ,5 ,67 ,117 ,114 ,118 ,101 ,1 ,16 ,0 ,1 ,1 ,88 ,1 ,255 ,138 ,0 ,1 ,1 ,89 ,1 ,255 ,138 ,0 ,0 ,0 ,10 ,255 ,137 ,5 ,1 ,2 ,255 ,140 ,0 ,0 ,0 ,82 ,255 ,130 ,1 ,10 ,112 ,97 ,114 ,101 ,110 ,116 ,72 ,97 ,115 ,104 ,2 ,1 ,25 ,99 ,114 ,121 ,112 ,116 ,111 ,47 ,101 ,108 ,108 ,105 ,112 ,116 ,105 ,99 ,46 ,112 ,50 ,50 ,52 ,67 ,117 ,114 ,118 ,101 ,255 ,141 ,3 ,1 ,1 ,9 ,112 ,50 ,50 ,52 ,67 ,117 ,114 ,118 ,101 ,1 ,255 ,142 ,0 ,1 ,1 ,1 ,11 ,67 ,117 ,114 ,118 ,101 ,80 ,97 ,114 ,97 ,109 ,115 ,1 ,255 ,144 ,0 ,0 ,0 ,83 ,255 ,143 ,3 ,1 ,1 ,11 ,67 ,117 ,114 ,118 ,101 ,80 ,97 ,114 ,97 ,109 ,115 ,1 ,255 ,144 ,0 ,1 ,7 ,1 ,1 ,80 ,1 ,255 ,138 ,0 ,1 ,1 ,78 ,1 ,255 ,138 ,0 ,1 ,1 ,66 ,1 ,255 ,138 ,0 ,1 ,2 ,71 ,120 ,1 ,255 ,138 ,0 ,1 ,2 ,71 ,121 ,1 ,255 ,138 ,0 ,1 ,7 ,66 ,105 ,116 ,83 ,105 ,122 ,101 ,1 ,4 ,0 ,1 ,4 ,78 ,97 ,109 ,101 ,1 ,12 ,0 ,0 ,0 ,255 ,239 ,255 ,142 ,255 ,169 ,1 ,1 ,29 ,2 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,1 ,1 ,29 ,2 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,255 ,22 ,162 ,224 ,184 ,240 ,62 ,19 ,221 ,41 ,69 ,92 ,92 ,42 ,61 ,1 ,29 ,2 ,180 ,5 ,10 ,133 ,12 ,4 ,179 ,171 ,245 ,65 ,50 ,86 ,80 ,68 ,176 ,183 ,215 ,191 ,216 ,186 ,39 ,11 ,57 ,67 ,35 ,85 ,255 ,180 ,1 ,29 ,2 ,183 ,14 ,12 ,189 ,107 ,180 ,191 ,127 ,50 ,19 ,144 ,185 ,74 ,3 ,193 ,211 ,86 ,194 ,17 ,34 ,52 ,50 ,128 ,214 ,17 ,92 ,29 ,33 ,1 ,29 ,2 ,189 ,55 ,99 ,136 ,181 ,247 ,35 ,251 ,76 ,34 ,223 ,230 ,205 ,67 ,117 ,160 ,90 ,7 ,71 ,100 ,68 ,213 ,129 ,153 ,133 ,0 ,126 ,52 ,1 ,254 ,1 ,192 ,1 ,5 ,80 ,45 ,50 ,50 ,52 ,0 ,0 ,1 ,29 ,2 ,65 ,148 ,109 ,36 ,229 ,243 ,53 ,205 ,245 ,246 ,83 ,142 ,214 ,61 ,73 ,151 ,143 ,232 ,199 ,172 ,102 ,59 ,239 ,66 ,45 ,233 ,118 ,219 ,1 ,29 ,2 ,48 ,46 ,153 ,110 ,55 ,45 ,171 ,244 ,39 ,227 ,168 ,174 ,65 ,214 ,53 ,236 ,232 ,77 ,71 ,7 ,137 ,132 ,229 ,19 ,204 ,34 ,194 ,64 ,0 ,1 ,4 ,0 ,4 ,0 ,0 ,0},
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
		if newBlock.MinerPublicKey != test.block.MinerPublicKey {
			t.Errorf("Bad Exit: Keys Don't match! '%v' vs. '%v'", newBlock.MinerPublicKey, test.block.MinerPublicKey)
		}
		if newBlock.Nonce != test.block.Nonce {
			t.Errorf("Bad Exit: Nonces Don't match! '%d' vs. '%d'", newBlock.Nonce, test.block.Nonce)
		}
	}
}
