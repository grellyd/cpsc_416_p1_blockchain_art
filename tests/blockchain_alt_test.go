package tests
/*
import (
	"fmt"
	"testing"
	"minerlib"
	"blockartlib"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func TestBlockChain (t *testing.T) {
	fmt.Println("Start test")

	privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	publicKey := &privateKey.PublicKey

	var set = blockartlib.MinerNetSettings{
		"83218ac34c1834c26781fe4bde918ee4",
		2,
		100,
		50,
		3000,
		5,
		5,
		blockartlib.CanvasSettings{0,0},
	}
	fmt.Println("set ", set)

	var canv = minerlib.CanvasData{}

	fmt.Println("canvas")
	var genBlock = minerlib.Block {
		"",
		nil,
		publicKey,
		0,
	}
	fmt.Println("genBlock")

	var miner = minerlib.Miner{
		0,
		nil,
		nil,
		nil,
		nil,
		publicKey,
		privateKey,
		minerlib.BCStorage{},
		&set,
		canv,
	}
	fmt.Println("miner")


	var bc = minerlib.NewBCRepresentation(&genBlock, &miner, "83218ac34c1834c26781fe4bde918ee4")
	fmt.Println("New Blockchain", bc)

	miner.Chain = *bc
	fmt.Println("miner: ", miner)

	var block1 = minerlib.Block{
		"83218ac34c1834c26781fe4bde918ee4",
		nil,
		publicKey,
		32,
	}

	var block2 = minerlib.Block{
		"83218ac34c1834c26781fe4bde900000",
		nil,
		publicKey,
		33,
	}

	var block3 = minerlib.Block{
		"83218ac34c1834c26781fe4bde900001",
		nil,
		publicKey,
		33,
	}

	var block4 = minerlib.Block{
		"83218ac34c1834c26781fe4bde900000",
		nil,
		publicKey,
		36,
	}

	var block5 = minerlib.Block{
		"83218ac34c1834c26781fe4bde900004",
		nil,
		publicKey,
		37,
	}

	var block6 = minerlib.Block{
		"83218ac34c1834c26781fe4bde900005",
		nil,
		publicKey,
		38,
	}

	// TESTS must create a single BC B1-B2-B3
	bc.AppendBlockToTree(&block1, &miner, "83218ac34c1834c26781fe4bde900000")
	fmt.Println("Tree1", bc.BCT)
	fmt.Println("Blockchain1", bc.BC.LastNode.Current.CurrentHash)

	bc.AppendBlockToTree(&block2, &miner, "83218ac34c1834c26781fe4bde900001")
	fmt.Println("Tree2", bc.BCT)
	fmt.Println("Blockchain2", bc.BC.LastNode.Current.CurrentHash)

	bc.AppendBlockToTree(&block3, &miner, "83218ac34c1834c26781fe4bde900002")
	fmt.Println("Tree3", bc.BCT)
	fmt.Println("Blockchain3", bc.BC.LastNode.Current.CurrentHash)

	// Now
	// B1 - B2 - B3
	//    - B4

	bc.AppendBlockToTree(&block4, &miner, "83218ac34c1834c26781fe4bde900004")
	fmt.Println("Tree4", bc.BCT)
	fmt.Println("Blockchain4", bc.BC.LastNode.Current.CurrentHash)

	// B1 - B2 - B3
	//    - B4 - B5
	b := bc.AppendBlockToTree(&block5, &miner, "83218ac34c1834c26781fe4bde900005")
	fmt.Println("Tree5", bc.BCT)
	fmt.Println("Blockchain4", bc.BC.LastNode.Current.CurrentHash)
	fmt.Println("BC shouldn't change: ", b)

	// B1 - B2 - B3
	//    - B4 - B5 - B6 <- longest
	b = bc.AppendBlockToTree(&block6, &miner, "83218ac34c1834c26781fe4bde900006")
	fmt.Println("Tree6", bc.BCT)
	fmt.Println("Blockchain6", bc.BC.LastNode.Current.CurrentHash)
	fmt.Println("BC should change: ", b)

	children, er := bc.GetChildrenNodes("83218ac34c1834c26781fe4bde900000")
	fmt.Println("Children ", children, "Error ", er)

	children, er = bc.GetChildrenNodes("83218ac34c1834c26781fe4bde900008")
	fmt.Println("Children ", children, "Error ", er)

	bc.AddToForest("83218ac34c1834c26781fe4bde900005", &block5)
	bc.AddToForest("83218ac34c1834c26781fe4bde900006", &block6)

	b = bc.IsInForest("83218ac34c1834c26781fe4bde900005")
	fmt.Println("Is in forest (true) ", b)

	bc.RemoveFromForest("83218ac34c1834c26781fe4bde900005")

	b = bc.IsInForest("83218ac34c1834c26781fe4bde900005")
	fmt.Println("Is in forest (false) ", b)


}*/