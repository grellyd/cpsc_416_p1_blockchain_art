package main

import (
	"blockartlib"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"keys"
	"networking"
	"minerlib"
	"net"
	"net/rpc"
	"os"
	"time"
)

var m minerlib.Miner // singleton for miner
var miners []net.Addr

var serverConnector *rpc.Client
var serverConn minerlib.ServerInstance
var artNodeConnector *rpc.Client
var OpQueue []*blockartlib.ArtNodeInstruction
var TreeQueue []string
	
func main() {
	gob.Register(elliptic.P384())
	fmt.Println("start")
	args := os.Args[1:]

	// Missing command line args.
	if len(args) != 3 {
		fmt.Println("usage: go run ink-miner.go [server ip:port] [pubKey] [privKey] ")
		return
	}
	/*
	Commented out to run locally. See Azure branch
	localIP := fmt.Sprintf("%s:8000", outboundIP)
	*/
	outboundIP :=  networking.GetOutboundIP()
	localIP := fmt.Sprintf("%s:0", outboundIP)
	keys, err := getKeyPair(args[2], args[1])
	CheckError(err)
	serverAddr, err := net.ResolveTCPAddr("tcp", args[0])
	CheckError(err)
	localAddr, err := net.ResolveTCPAddr("tcp", localIP)
	CheckError(err)
	
	// Create Miner
	m = minerlib.NewMiner(serverAddr, keys)
	
	//setup an ArtNode Reciever
	artNodeInst := new(ArtNodeInstance)
	// register art node instance locally
	rpc.Register(artNodeInst)

	// Add a listener on myself
	localListener, err := net.ListenTCP("tcp", localAddr)
	CheckError(err)
	fmt.Println("Local addr: ", localListener.Addr().String())

	// My Info to Send
	localMinerInfo := MinerInfo{localListener.Addr(), m.PublKey}
	m.ServerNodeAddr, _ = net.ResolveTCPAddr("tcp", localMinerInfo.Address.String())
	fmt.Println("Serv addr: ", m.ServerNodeAddr.String())
	
	// Connect to server
	serverConn, err = connectServer(serverAddr, localMinerInfo, m.Settings)
	CheckError(err)
	fmt.Println("Settings ", m.Settings)

	// Setup Heartbeats
	go doEvery(time.Duration(m.Settings.HeartBeat-1000)*time.Millisecond, serverConn.SendHeartbeat)

	// TODO: Check in with neighbours
	// TODO: Ask Neighbours for blockchain that already exists
	genesisBlock := m.CreateGenesisBlock()

	m.Blockchain = minerlib.NewBlockchainStorage(genesisBlock, m.Settings)
	CheckError(err)
	go m.StartMining()
	//go m.TestEarlyExit()

	// Ask for Neighbours
	fmt.Println("Asking for neighbours")
	err = serverConn.RequestMiners(&miners, m.Settings.MinNumMinerConnections)
	fmt.Println("Got neighbours!")
	CheckError(err)
	fmt.Println("miners1: ", miners)

	err = m.AddMinersToList(&miners)
	CheckError(err)
	fmt.Printf("miners1: %v \n", &m.Neighbours)

	if len(m.Neighbours) !=0 {
		err = m.OpenNeighbourConnections()
		CheckError(err)
		fmt.Println("Opened RPC connections to neighbour miners")

		neighbourToReceiveBCFrom, err := m.ConnectToNeighbourMiners(m.ServerNodeAddr)
		CheckError(err)
		fmt.Printf("Connected to neighbour miners; will ask for BlockChain from neighbour with address %s\n", neighbourToReceiveBCFrom.String())

		err = m.RequestBCStorageFromNeighbour(&neighbourToReceiveBCFrom, &TreeQueue)
		CheckError(err)
		fmt.Println("Requested BCStorage from neighbour")
	}

	// Set up receiver for other Miners
	minerReceiverInst := new(MinerInstance)
	rpc.Register(minerReceiverInst)

	fmt.Printf("befor goRoutine: %v aaaand length %v, \n", &m.Neighbours, len(m.Neighbours))
	go doEvery(5*time.Second, UpdateNeighbours)

	serviceRequests(localListener)
}

func connectServer(serverAddr *net.TCPAddr, minerInfo MinerInfo, settings *blockartlib.MinerNetSettings) (serverConnection minerlib.ServerInstance, err error) {
	// dial to server
	serverRPCClient, err := rpc.Dial("tcp", serverAddr.String())
	CheckError(err)
	// setup gob
	gob.Register(&net.TCPAddr{})
	gob.Register(&elliptic.CurveParams{})

	//1st rpc call
	//2nd retrieve settings ==> 2 in 1
	err = serverRPCClient.Call("RServer.Register", minerInfo, settings)
	CheckError(err)
	// Create the serverConnection. 
	// TODO: refactor to ServerInstance
	tcpFromAddr, err := net.ResolveTCPAddr("tcp", minerInfo.Address.String())
	CheckError(err)
	serverConnection = minerlib.ServerInstance{
		Addr:      *tcpFromAddr,
		RPCClient: serverRPCClient,
		Public:    minerInfo.Key,
	}
	return serverConnection, nil
}

func serviceRequests(localListener *net.TCPListener) {
	for {
		conn, err := localListener.Accept()
		CheckError(err)
		defer conn.Close()

		go rpc.ServeConn(conn)
		fmt.Println("after connection served")

		time.Sleep(10 * time.Millisecond)

		if len(OpQueue) != 0 {
			fmt.Println("connect to queue")
			artNodeConnector, err = rpc.Dial("tcp", OpQueue[0].LocalIP)
			CheckError(err)
			var b bool
			err = artNodeConnector.Call("MinerInstance.ConnectMiner", true, &b)
			CheckError(err)
			OpQueue = OpQueue[1:]
			fmt.Println("connected to queue ", b, "len ", len(OpQueue))
		}
		// DrawOperations to validate
		// For valid add to miner op channel
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func doEvery(d time.Duration, f func(time.Time) error) error {
	for x := range time.Tick(d) {
		f(x)
	}
	return nil
}

func getKeyPair(privStr string, pubStr string) (*blockartlib.KeyPair, error) {
	priv, pub := keys.Decode(privStr, pubStr)
	pair := blockartlib.KeyPair{
		Private: priv,
		Public:  pub,
	}
	return &pair, nil

}

func UpdateNeighbours(t time.Time) (err error) {
	lom := make([]net.Addr, 0)
	fmt.Printf("start updateN, lom %v lenLom %v, minersN %v \n ", &lom, len(lom), len(m.Neighbours))

	if len(m.Neighbours) < int(m.Settings.MinNumMinerConnections) {
		fmt.Println("starting request again")
		err = serverConn.RequestMiners(&lom, m.Settings.MinNumMinerConnections)
		fmt.Printf("starting request again, lom %v lenLom %v, minersN %v \n ", &lom, len(lom), len(m.Neighbours))
		if len(lom) !=0 {m.AddMinersToList(&lom)} else {
			return nil
		}
	}
	if len(m.Neighbours) == 0 {
		return nil
	}
	if allAlive(&m) {
		return nil
	}
	fmt.Printf("Neigh addr: %v \n", &m.Neighbours)
	e := m.OpenNeighbourConnections()
	CheckError(e)
	fmt.Println("Server node address: ", m.ServerNodeAddr.String())
	neighbourToReceiveBCFrom, err := m.ConnectToNeighbourMiners(m.ServerNodeAddr)
	CheckError(err)
	fmt.Printf("Connected to neighbour miners in Update; will ask for BlockChain from neighbour with address %s\n", neighbourToReceiveBCFrom.String())

	if neighbourToReceiveBCFrom.Port == 0 {
		return nil
	}
	err = m.RequestBCStorageFromNeighbour(&neighbourToReceiveBCFrom, &TreeQueue)
	CheckError(err)
	fmt.Println("Requested BCStorage from neighbour in Update")

	if err != nil {
		return err
	}
	return nil
}

func allAlive(m *minerlib.Miner) bool {
	for _,v := range m.Neighbours {
		if !v.Alive {return false}
	}
	return true
}

// =========================
// Connection Instances
// =========================

// RPC Connections with ArtNode
type ArtNodeInstance int // same as above

func (si *ArtNodeInstance) ConnectNode(an *blockartlib.ArtNodeInstruction, reply *bool) error {
	fmt.Println("In rpc call to register the AN")
	privateKey := keys.DecodePrivateKey(an.PrivKey)
	publicKey := keys.DecodePublicKey(an.PubKey)
	// TODO check if already connected
	if !keys.MatchingPair(privateKey, publicKey) {
		fmt.Println("Invalid key pair.")
		return blockartlib.DisconnectedError("Key pair isn't valid")
	}else {
		*reply = true
		OpQueue = append(OpQueue, an)
	}
	return nil
}

func (si *ArtNodeInstance) GetGenesisBlockHash(stub *bool, reply *string) error {
	fmt.Println("In RPC getting hash of genesis block")
	// TODO: check if connected
	*reply = m.Settings.GenesisBlockHash
	return nil
}

func (si *ArtNodeInstance) GetAvailableInk(stub *bool, reply *uint32) error {
	fmt.Println("In RPC getting ink from miner")
	*reply = m.InkLevel
	return nil
}

func (si *ArtNodeInstance) GetSVGString(shapeHash string, reply *string) error {
	fmt.Println("In RPC getting svg string")
	//m.Blockchain.BC
	temp := m.Blockchain.BC.GenesisNode
	var b *minerlib.Block
	for {
		if temp.Current == nil {
			break
		}
		b = temp.Current.BlockResiding
		for _, op := range b.Operations {
			if op.ShapeHash == shapeHash {
				//*reply = minerlib.OpToSvg(*op, m.Settings.CanvasSettings)
				*reply = "<path d=\"M5,2L6,2L6,3L5,3L5,2\" stroke=\"red\" fill =\"transparent\"/>"
				return nil
			}
		}
	}
	return blockartlib.InvalidShapeHashError(shapeHash)
}

func (si *ArtNodeInstance) GetAllSVGStrings(blockHash string, reply []string) error {
	fmt.Println("In RPC getting svg string")
	//treeNode := minerlib.FindBCTreeNode(m.Blockchain.BCT.GenesisNode, *blockHash)

	// iterate over blockchain, get all svg strings
	reply = append(reply, "<path d=\"M5,2L6,2L6,3L5,3L5,2\" stroke=\"red\" fill =\"transparent\"/>")
	return err
}

func (si *ArtNodeInstance) GetBlockChildren(hash *string, reply *[]string) error {
	fmt.Println("In RPC getting children hashes")
	bla, err := m.Blockchain.GetChildrenNodes(*hash)
	*reply = bla
	CheckError(err)
	return err
}

func (si *ArtNodeInstance) SubmitOperation(op blockartlib.Operation, shapeHash *string) error {
	// TODO use the an connection with the channel to wait
	err := m.AddOp(&op)
	if err != nil {
		return fmt.Errorf("unable to submit operation: %v", err)
	}
	// blocks until done at validation depth
	hash, err := m.GetShapeHash(&op)
	shapeHash = &hash
	return err
}

func (si *ArtNodeInstance) GetShapesFromBlock (blockHash *string, reply *[]string) error {
	fmt.Println("In RPC getting shape from block")
	treeNode := minerlib.FindBCTreeNode(m.Blockchain.BCT.GenesisNode, *blockHash)
	if treeNode == nil {
		return blockartlib.InvalidBlockHashError("invalid hash")
	}
	block := treeNode.BlockResiding
	ops := block.Operations
	for _,v := range ops {
		*reply = append(*reply, v.ShapeHash)
	}
	return nil
}

// RPC Connections with other Miners
type MinerInstance int

func (si *MinerInstance) ConnectNewNeighbour(neighbourAddr *net.TCPAddr, reply *int) error {
	// Add neighbour to list of neighbours
	fmt.Printf("Got request to register a new neighbour with TCP address %s\n", neighbourAddr.String())
	/*var newNeighbour = minerlib.MinerConnection{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", neighbourAddr.String())
	if err != nil {
		return err
	}*/
	//newNeighbour.Addr = *tcpAddr
	//m.Neighbours = append(m.Neighbours, &newNeighbour)
	lom := make([]net.Addr, 0)
	lom = append(lom, neighbourAddr)
	e := m.AddMinersToList(&lom)
	CheckError(e)

	// Return the length of our blockchain, so the new neighbour can decide
	// if they want our tree.

	*reply = m.Blockchain.BC.LastNode.Current.Depth
	fmt.Printf("ConnectNewNeighbour: Returning a reply depth of %d\n", *reply)

	return nil
}

// TODO: check switch
func (mi *MinerInstance) ReceiveBlockFromNeighbour (blockMarshalled *[][]byte, reply *bool) error {
	block, err := minerlib.UnmarshallBinary(*blockMarshalled)
	CheckError(err)
	m.AddDisseminatedBlock(block)
	*reply = true
	return err
}

func (mi *MinerInstance) ReceiveOpFromNeighbour(opMarshalled *[]byte, reply *bool) error {
	_, err := blockartlib.OperationUnmarshall(*opMarshalled)
	CheckError(err)
	return err
}

func (mi *MinerInstance) DisseminateTree (variable *bool, reply *[]string) error {
	*reply = m.MarshallTree(reply, nil)
	return nil
}

func (mi *MinerInstance) GiveBlock (blockHash *string, reply *[][]byte) error {
	block, err := m.Blockchain.FindBlockByHash(*blockHash)
	if err !=nil {
		fmt.Println("Error in GiveBlock RPC: ", err)
	}
	*reply, err = block.MarshallBinary()
	return err
}

// struct for communicating info about a miner to the server
type MinerInfo struct {
	Address net.Addr
	Key     *ecdsa.PublicKey
}

