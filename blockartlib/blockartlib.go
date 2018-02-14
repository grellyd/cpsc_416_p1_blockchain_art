/*

This package specifies the application's interface to the the BlockArt
library (blockartlib) to be used in project 1 of UBC CS 416 2017W2.

*/

package blockartlib

import (
	"crypto/ecdsa"
	"time"
)

// Represents a type of shape in the BlockArt system.
// 1 for a line, 2 for a circle
type ShapeType int

const (
	// Path shape.
	PATH ShapeType = iota

	// Circle shape (extra credit).
	CIRCLE
)

// Settings for a canvas in BlockArt.
type CanvasSettings struct {
	// Canvas dimensions
	CanvasXMax uint32
	CanvasYMax uint32
}

type KeyPair struct {
	Private *ecdsa.PrivateKey
	Public *ecdsa.PublicKey
}

type OperationType int

const (
	NOP OperationType = iota
	DRAW
	DELETE
)

type Operation struct {
	Type OperationType
	OperationNumber int
	OperationSig string
	Shape ShapeType
	Fill string // Can be "transparent" or "filled"
	Stroke string
	ShapeSVGString string
	ArtNodePubKey string
	Nonce uint32
}

func (o *Operation) GenerateSig() error {
	// TODO: A function of the Op and the pub key
	return nil
}

type ArtNodeInstruction struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
	LocalIP		string
}

type MinerInstructions struct {
	Operation Operation
}

// Settings for an instance of the BlockArt project/network.
type MinerNetSettings struct {
	// Hash of the very first (empty) block in the chain.
	GenesisBlockHash string

	// The minimum number of ink miners that an ink miner should be
	// connected to. If the ink miner dips below this number, then
	// they have to retrieve more nodes from the server using
	// GetNodes().
	MinNumMinerConnections uint8

	// Mining ink reward per op and no-op blocks (>= 1)
	InkPerOpBlock   uint32
	InkPerNoOpBlock uint32

	// Number of milliseconds between heartbeat messages to the server.
	HeartBeat uint32

	// Proof of work difficulty: number of zeroes in prefix (>=0)
	PoWDifficultyOpBlock   uint8
	PoWDifficultyNoOpBlock uint8

	// Canvas settings
	CanvasSettings CanvasSettings
}



// Represents a canvas in the system.
type Canvas interface {
	// Adds a new shape to the canvas.
	// Can return the following errors:
	// - DisconnectedError
	// - InsufficientInkError
	// - InvalidShapeSvgStringError
	// - ShapeSvgStringTooLongError
	AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error)

	// Returns the encoding of the shape as an svg string.
	// Can return the following errors:
	// - DisconnectedError
	// - InvalidShapeHashError
	GetSvgString(shapeHash string) (svgString string, err error)

	// Returns the amount of ink currently available.
	// Can return the following errors:
	// - DisconnectedError
	GetInk() (inkRemaining uint32, err error)

	// Removes a shape from the canvas.
	// Can return the following errors:
	// - DisconnectedError
	// - ShapeOwnerError
	// - OutOfBoundsError
	// - ShapeOverlapError
	DeleteShape(validateNum uint8, shapeHash string) (inkRemaining uint32, err error)

	// Retrieves hashes contained by a specific block.
	// Can return the following errors:
	// - DisconnectedError
	// - InvalidBlockHashError
	GetShapes(blockHash string) (shapeHashes []string, err error)

	// Returns the block hash of the genesis block.
	// Can return the following errors:
	// - DisconnectedError
	GetGenesisBlock() (blockHash string, err error)

	// Retrieves the children blocks of the block identified by blockHash.
	// Can return the following errors:
	// - DisconnectedError
	// - InvalidBlockHashError
	GetChildren(blockHash string) (blockHashes []string, err error)

	// Closes the canvas/connection to the BlockArt network.
	// - DisconnectedError
	CloseCanvas() (inkRemaining uint32, err error)
}

// The constructor for a new Canvas object instance. Takes the miner's
// IP:port address string and a public-private key pair (ecdsa private
// key type contains the public key). Returns a Canvas instance that
// can be used for all future interactions with blockartlib.
//
// The returned Canvas instance is a singleton: an application is
// expected to interact with just one Canvas instance at a time.
//
// Can return the following errors:
// - DisconnectedError
func OpenCanvas(minerAddr string, privKey ecdsa.PrivateKey) (canvas Canvas, setting CanvasSettings, err error) {
	// TODO
	// For now return DisconnectedError
	println("Started canvas")

	var an = ArtNode {
		0,
		minerAddr,
		&privKey,
		&privKey.PublicKey,
		false,
		"127.0.0.1" + ":0",
		nil}
	//var settings CanvasSettings
	err = an.Connect(an.MinerAddr, an.PrivKey)
	CheckErr(err)

	/*err = an.MinerConnection.Call("ArtNodeInstance.GetCanvasSetting", &an, &settings)
	CheckErr(err)*/

	println("Miner addr ", an.MinerAddr, "LocalIP: ", an.LocalIP, "Miner Connection ", an.MinerConnection, "Alive: ", an.MinerAlive, "Error ", err )
	time.Sleep(5*time.Second)
	if err == nil {
		return an, CanvasSettings{}, nil
	}
	return nil, CanvasSettings{}, DisconnectedError("")
}


