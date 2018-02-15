package blockartlib

import (
	"keys"
	"net/rpc"
	"net"
	"fmt"
	"crypto/ecdsa"
	"crypto/x509"
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

// RPC calls
func (mi *MinerInstance) ConnectMiner (mins *bool, reply *bool) error {
	fmt.Println("In RPC connecting Miner")
	*reply = true
	return nil
}

// CArtNodeVAS INTERFACE FUNCTIONS
func (an ArtNode) AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error) {
	op := Operation{
		Type: DRAW,
		OperationNumber: 0,
		Shape: shapeType,
		Fill: fill,
		Stroke: stroke,
		ShapeSVGString: shapeSvgString,
		ArtNodePubKey: keys.EncodePublicKey(an.PubKey),
		ValidateBlockNum: validateNum,
		ShapeHash: "",
	}
	op.GetNumber()
	op.GenerateSig()
	err = an.MinerConnection.Call("ArtNodeInstance.SubmitOperation", op, &blockHash)
	if err != nil {
		return "", "", 0, fmt.Errorf("unable to submit operation: %v", err)
	}
	inkRemaining, err = an.GetInk()
	if err != nil {
		return "", "", 0, fmt.Errorf("unable to get ink: %v", err)
	}
	shapeHashes, err := an.GetShapes(blockHash)
	if len(shapeHashes) == 0 || err != nil {
		return "", "", 0, fmt.Errorf("no shapehashes returned: %v", err)
	}
	// TODO validate there is only one and pick if needed
	shapeHash = shapeHashes[0]
	return shapeHash, blockHash, inkRemaining, nil
}

func (an ArtNode) GetSvgString(shapeHash string) (svgString string, err error) {
	return svgString, err
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
		// TODO set operation number
		OperationNumber: 0,
		Shape: 0,
		Fill: "",
		Stroke: "",
		ShapeSVGString: "",
		ArtNodePubKey: keys.EncodePublicKey(an.PubKey),
		ValidateBlockNum: validateNum,
		ShapeHash: shapeHash,
	}
	err = op.GetNumber()
	if err != nil {
		return 0, fmt.Errorf("unable to set operation number: %v", err)
	}
	err = op.GenerateSig()
	if err != nil {
		return 0, fmt.Errorf("unable to generate operation sig: %v", err)
	}

	blockHash := ""
	err = an.MinerConnection.Call("ArtNodeInstance.SubmitOperation", op, &blockHash)
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
	return inkRemaining, err
}

// MINER INTERACTION FUNCTIONS
func (an *ArtNode) Connect(minerAddr string, privKey *ecdsa.PrivateKey) (err error) {
	// Establish RPC connection
	minerInst := new(MinerInstance)
	rpc.Register(minerInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", an.LocalIP)
	CheckErr(err)
	fmt.Println("TCP: ", tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckErr(err)
	fmt.Println("listening on", listener.Addr().String())
	an.LocalIP = listener.Addr().String()

	go rpc.Accept(listener)

	// connect to the miner
	an.MinerConnection, err = rpc.Dial("tcp", an.MinerAddr)
	CheckErr(err)


	fmt.Println("Miner Connection: ", an.MinerConnection)
	var reply bool // TODO: change when actual RPC will be alive
	//gob.RegisterName("crypto/elliptic.CurveParams", elliptic.CurveParams{})
	//gob.Register(elliptic.CurveParams{})

	pk, _ := x509.MarshalECPrivateKey(an.PrivKey)
	puk, _ := x509.MarshalPKIXPublicKey(an.PubKey)
	an1 := ArtNodeInstruction{
		0,
		an.MinerAddr,
		string(pk),
		string(puk),
		false,
		an.LocalIP,
	}
	fmt.Println("trying to connect via rpc")
	err = an.MinerConnection.Call("ArtNodeInstance.ConnectNode", an1, &reply)
	CheckErr(err)
	fmt.Println("connected via rpc ", reply)
	an.MinerAlive = true
	//time.Sleep(1*time.Second)
	return nil
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

