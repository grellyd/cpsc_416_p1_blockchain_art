package blockartlib

import (
	"net/rpc"
	"net"
	"fmt"
)

/*
Artnode that communicates with the client app and the miner
Implements the canvas interface
In short, the work for the client is in here
*/

type MinerInst struct {
}

type ArtNode struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
	LocalIP		string
	MinerConnection *rpc.Client
}


// CArtNodeVAS INTERFACE FUNCTIONS
func (an *ArtNode) AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error) {
	return shapeHash, blockHash, inkRemaining, err
}

func (an *ArtNode) GetSvgString(shapeHash string) (svgString string, err error) {
	return svgString, err
}

func (an *ArtNode) GetInk() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *ArtNode) DeleteShape(validateNum uint8, shapeHash string) (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *ArtNode) GetShapes(blockHash string) (shapeHashes []string, err error) {
	return shapeHashes, err
}

func (an *ArtNode) GetGenesisBlock() (blockHash string, err error) {
	return blockHash, err
}

func (an *ArtNode) GetChildren(blockHash string) (blockHashes []string, err error) {
	return blockHashes, err
}
func (an *ArtNode) CloseCanvas() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

// MINER INTERACTION FUNCTIONS
func (an *ArtNode) Connect(minerAddr, pubKey, privKey string) (err error) {
	// Establish RPC connection
	minerInst := new(MinerInst)
	rpc.Register(minerInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", an.LocalIP)
	CheckErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckErr(err)

	go rpc.Accept(listener)

	// connect to the miner
	an.MinerConnection, err = rpc.Dial("tcp", an.MinerAddr)
	CheckErr(err)

	var reply bool // TODO: change when actual RPC will be alive
	err = an.MinerConnection.Call("ArtNodeInst.ConnectNode", an, &reply)

	return nil
}

func (an *ArtNode) MakeDrawRequest() (err error) {
	return err
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

