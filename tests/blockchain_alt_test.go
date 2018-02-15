package tests
import (
	"fmt"
	"testing"
	"minerlib"
	"blockartlib"
	"keys"
)

func TestBlockChain (t *testing.T) {
	fmt.Println("Start test")

	privateKey, publicKey, _ := keys.Generate()

	var set = blockartlib.MinerNetSettings{
		GenesisBlockHash: "83218ac34c1834c26781fe4bde918ee4",
		MinNumMinerConnections: 2,
		InkPerOpBlock: 100,
		InkPerNoOpBlock: 50,
		HeartBeat: 3000,
		PoWDifficultyOpBlock: 5,
		PoWDifficultyNoOpBlock: 5,
		CanvasSettings: blockartlib.CanvasSettings{
			CanvasXMax: 0,
			CanvasYMax: 0,
		},
	}
	fmt.Println("set ", set)

	var canv = minerlib.CanvasData{}

	fmt.Println("canvas")

	var miner = minerlib.Miner{
		InkLevel: 0,
		ServerNodeAddr: nil,
		ServerHrtBtAddr: nil,
		ArtNodes: nil,
		Neighbors: nil,
		PublKey: publicKey,
		PrivKey: privateKey,
		Blockchain: &minerlib.BCStorage{},
		Settings: &set,
		LocalCanvas: canv,
	}
	fmt.Println("miner")
	genBlock := miner.CreateGenesisBlock()
	fmt.Println("genBlock")

	// The following hashes have to change
	var bc = minerlib.NewBlockchainStorage(genBlock, &set)
	fmt.Println("New Blockchain", bc)

	miner.Blockchain = bc
	fmt.Println("miner: ", miner)

	var block1 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde918ee4",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 32,
	}

	var block2 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde900000",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 33,
	}

	var block3 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde900001",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 33,
	}

	var block4 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde900000",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 36,
	}

	var block5 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde900004",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 37,
	}

	var block6 = minerlib.Block{
		ParentHash: "83218ac34c1834c26781fe4bde900005",
		Operations: nil,
		MinerPublicKey: publicKey,
		Nonce: 38,
	}

	// TESTS must create a single BC B1-B2-B3
	res := bc.AppendBlock(&block1, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block1, bc)
	}
	fmt.Println("Tree1", bc.BCT)
	fmt.Println("Blockchain1", bc.BC.LastNode.Current.CurrentHash)

	res = bc.AppendBlock(&block2, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block2, bc)
	}
	fmt.Println("Tree2", bc.BCT)
	fmt.Println("Blockchain2", bc.BC.LastNode.Current.CurrentHash)

	res = bc.AppendBlock(&block3, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block3, bc)
	}
	fmt.Println("Tree3", bc.BCT)
	fmt.Println("Blockchain3", bc.BC.LastNode.Current.CurrentHash)

	// Now
	// B1 - B2 - B3
	//    - B4

	res = bc.AppendBlock(&block4, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block4, bc)
	}
	fmt.Println("Tree4", bc.BCT)
	fmt.Println("Blockchain4", bc.BC.LastNode.Current.CurrentHash)

	// B1 - B2 - B3
	//    - B4 - B5
	res = bc.AppendBlock(&block5, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block5, bc)
	}
	fmt.Println("Tree5", bc.BCT)
	fmt.Println("Blockchain4", bc.BC.LastNode.Current.CurrentHash)
	// fmt.Println("BC shouldn't change: ", b)

	// B1 - B2 - B3
	//    - B4 - B5 - B6 <- longest
	res = bc.AppendBlock(&block6, miner.Settings)
	if !res {
		t.Errorf("Error when appending block '%v' to '%v'", block6, bc)
	}
	fmt.Println("Tree6", bc.BCT)
	fmt.Println("Blockchain6", bc.BC.LastNode.Current.CurrentHash)
	// fmt.Println("BC should change: ", b)

	children, er := bc.GetChildrenNodes("83218ac34c1834c26781fe4bde900000")
	fmt.Println("Children ", children, "Error ", er)

	children, er = bc.GetChildrenNodes("83218ac34c1834c26781fe4bde900008")
	fmt.Println("Children ", children, "Error ", er)

	/*
	bc.AddToForest("83218ac34c1834c26781fe4bde900005", &block5)
	bc.AddToForest("83218ac34c1834c26781fe4bde900006", &block6)

	b = bc.IsInForest("83218ac34c1834c26781fe4bde900005")
	fmt.Println("Is in forest (true) ", b)

	bc.RemoveFromForest("83218ac34c1834c26781fe4bde900005")

	b = bc.IsInForest("83218ac34c1834c26781fe4bde900005")
	fmt.Println("Is in forest (false) ", b)
	*/
}
