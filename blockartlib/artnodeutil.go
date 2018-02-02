package blockartlib

import (
	"net/rpc"
	"net"
	"fmt"
)

/*
Artnode that communicates with the client app and the miner
*/

type AN struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
	LocalIP		string
}

type MinerInst int // for now it's the int, but we can change to actual struct
var minerConnector *rpc.Client

// CANVAS INTERFACE FUNCTIONS
func (an *AN) AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error) {
	return shapeHash, blockHash, inkRemaining, err
}

func (an *AN) GetSvgString(shapeHash string) (svgString string, err error) {
	return svgString, err
}

func (an *AN) GetInk() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *AN) DeleteShape(validateNum uint8, shapeHash string) (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *AN) GetShapes(blockHash string) (shapeHashes []string, err error) {
	return shapeHashes, err
}

func (an *AN) GetGenesisBlock() (blockHash string, err error) {
	return blockHash, err
}

func (an *AN) GetChildren(blockHash string) (blockHashes []string, err error) {
	return blockHashes, err
}
func (an *AN) CloseCanvas() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

// MINER INTERACTION FUNCTIONS
func (an *AN) Connect(minerAddr, pubKey, privKey string) (err error) {
	// Establish RPC connection
	minerInst := new(MinerInst)
	rpc.Register(minerInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", an.LocalIP)
	CheckErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckErr(err)

	go rpc.Accept(listener)

	// connect to the miner
	minerConnector, err = rpc.Dial("tcp", an.MinerAddr)
	CheckErr(err)

	var reply bool // TODO: change when actual RPC will be alive
	err = minerConnector.Call("ArtNodeInst.ConnectNode", an, &reply)


	return nil
}

func (an *AN) MakeDrawRequest() (err error) {
	return err
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

