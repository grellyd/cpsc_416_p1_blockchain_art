package blockartlib

import (
	"os"
	"keys"
	"net/rpc"
	"net"
	"fmt"
	"crypto/ecdsa"
)

/*
Artnode that communicates with the client app and the miner
Implements the canvas interface
In short, the work for the client is in here
*/

type MinerInstance int

type ArtNode struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	*ecdsa.PrivateKey
	PubKey 		*ecdsa.PublicKey
	MinerAlive 	bool
	LocalIP		string
	MinerConnection *rpc.Client
}

// eww global but the Art Node is pass by value, not by reference; changes in a function don't persist
var opNumber uint32 = 0

// RPC calls
func (mi *MinerInstance) ConnectMiner (mins *bool, reply *bool) error {
	fmt.Println("In RPC connecting Miner")
	*reply = true
	return nil
}

// CArtNodeVAS INTERFACE FUNCTIONS
func (an ArtNode) AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error) {
	fmt.Println("ARTNODEUTIL: Calling AddShape")
	op := Operation{
		Type: DRAW,
		OperationNumber: opNumber,
		Shape: shapeType,
		Fill: fill,
		Stroke: stroke,
		ShapeSVGString: shapeSvgString,
		ArtNodePubKey: keys.EncodePublicKey(an.PubKey),
		ValidateBlockNum: validateNum,
		ShapeHash: "",
	}
	opNumber = opNumber + 1
	err = op.GenerateSig()
	if err != nil {
		return "", "", 0, fmt.Errorf("unable to generate operation sig: %v", err)
	}
	fmt.Printf("[artnodeutil] op: %v\n", op)
	fmt.Println("ARTNODEUTIL: Calling RPC call to Miner: ArtNodeInstance.SubmitOperation")
	err = an.MinerConnection.Call("ArtNodeInstance.SubmitOperation", op, &shapeHash)
	if err != nil {
		return "", "", 0, fmt.Errorf("ARTNODEUTIL.AddShape: unable to submit operation: %v", err)
	}
	fmt.Printf("[artnodeutil] shapeHash: '%v'\n", shapeHash)
	inkRemaining, err = an.GetInk()
	if err != nil {
		return "", "", 0, fmt.Errorf("unable to get ink: %v", err)
	}

	return shapeHash, blockHash, inkRemaining, nil
}

func (an ArtNode) GetSvgString(shapeHash string) (svgString string, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetSVGString", shapeHash, &svgString)
	if err != nil {
		return "", DisconnectedError("miner unavailable")
	}
	return svgString, err
}

func (an ArtNode) GetAllSvgStrings(blockHash string) (svgStrings []string, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetAllSVGStrings", blockHash, &svgStrings)
	if err != nil {
		return svgStrings, DisconnectedError("miner unavailable")
	}
	return svgStrings, err
}

func (an ArtNode) GetInk() (inkRemaining uint32, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetAvailableInk", true, &inkRemaining)
	if err != nil {
		return 0, DisconnectedError("miner unavailable")
	}
	return inkRemaining, err
}

func (an ArtNode) DeleteShape(validateNum uint8, shapeHash string) (inkRemaining uint32, err error) {
	op := Operation{
		Type: DELETE,
		OperationNumber: opNumber,
		Shape: 0,
		Fill: "",
		Stroke: "",
		ShapeSVGString: "",
		ArtNodePubKey: keys.EncodePublicKey(an.PubKey),
		ValidateBlockNum: validateNum,
		ShapeHash: shapeHash,
	}
	opNumber = opNumber + 1

	fmt.Printf("[artnodeutil#DeleteShape] ShapeHash: '%v'\n", op.ShapeHash)
	err = an.MinerConnection.Call("ArtNodeInstance.SubmitOperation", op, &shapeHash)
	if err != nil {
		return 0, fmt.Errorf("unable to submit operation: %v", err)
	}

	inkRemaining, err = an.GetInk()
	if err != nil {
		return 0, fmt.Errorf("unable to get ink: %v", err)
	}
	return inkRemaining, err
}

func (an ArtNode) GetShapes(blockHash string) (shapeHashes []string, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetShapesFromBlock", &blockHash, &shapeHashes)
	if err != nil {
		return shapeHashes, DisconnectedError("miner unavailable") // TODO: check type of error
	}
	return shapeHashes, err
}

func (an ArtNode) GetGenesisBlock() (blockHash string, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetGenesisBlockHash", true, &blockHash)
	if err != nil {
		return "", DisconnectedError("miner unavailable")
	}
	return blockHash, err
}

func (an ArtNode) GetChildren(blockHash string) (blockHashes []string, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetBlockChildren", blockHash, &blockHashes)
	if err != nil {
		return blockHashes, DisconnectedError("miner unavailable")
	}
	return blockHashes, err
}

func (an ArtNode) CloseCanvas() (inkRemaining uint32, err error) {
	err = an.MinerConnection.Call("ArtNodeInstance.GetAvailableInk", true, &inkRemaining)
	if err != nil {
		return 0, DisconnectedError("miner unavailable")
	}
	an.MinerConnection.Close()
	return inkRemaining, err
}

// MINER INTERACTION FUNCTIONS
func (an *ArtNode) Connect(minerAddr string, privKey *ecdsa.PrivateKey) (err error) {
	fmt.Println("ARTNODEUTIL: Running Connect to connect to miner at address ", minerAddr)
	// Establish RPC connection to Miner
	minerInst := new(MinerInstance)
	rpc.Register(minerInst)

	fmt.Printf("ARTNODEUTIL: Resolving ArtNode local/outbound IP: %s\n", an.LocalIP)
	tcpAddr, err := net.ResolveTCPAddr("tcp", an.LocalIP)
	CheckErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckErr(err)
	fmt.Println("ARTNODEUTIL: ArtNode listening on", listener.Addr().String())

	go rpc.Accept(listener)

	// connect to the miner
	an.MinerConnection, err = rpc.Dial("tcp", an.MinerAddr)

	CheckErr(err)

	var reply bool // TODO: change when actual RPC will be alive
	//gob.RegisterName("crypto/elliptic.CurveParams", elliptic.CurveParams{})
	//gob.Register(elliptic.CurveParams{})

	pk := keys.EncodePrivateKey(an.PrivKey)
	puk := keys.EncodePublicKey(an.PubKey)
	an1 := ArtNodeInstruction{
		0,
		an.MinerAddr,
		string(pk),
		string(puk),
		false,
		listener.Addr().String(),
	}
	fmt.Println("ARTNODEUTIL: trying to connec to Miner via RPC: ArtNodeInstance.ConnectNode")
	err = an.MinerConnection.Call("ArtNodeInstance.ConnectNode", an1, &reply)
	CheckErr(err)
	fmt.Println("ARTNODEUTIL connected via rpc without error; reply is: ", reply)
	an.MinerAlive = true
	//time.Sleep(1*time.Second)
	return nil
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

