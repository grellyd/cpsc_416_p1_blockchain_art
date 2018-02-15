package main

import (
	"blockartlib"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"keys"
	"minerlib"
	"net"
	"net/rpc"
	"os"
	"time"
)

var m minerlib.Miner // singleton for miner
var miners []net.Addr

var serverConnector *rpc.Client
var artNodeConnector *rpc.Client
var OpQueue []*blockartlib.ArtNodeInstruction

var localIP = "127.0.0.1:0"

func main() {
	fmt.Println("start")
	args := os.Args[1:]

	// Missing command line args.
	if len(args) != 3 {
		fmt.Println("usage: go run ink-miner.go [server ip:port] [pubKey] [privKey] ")
		return
	}
	keys, err := getKeyPair(args[2], args[1])
	checkError(err)
	serverAddr, err := net.ResolveTCPAddr("tcp", args[0])
	checkError(err)
	localAddr, err := net.ResolveTCPAddr("tcp", localIP)
	checkError(err)

	// Create Miner
	m = minerlib.NewMiner(serverAddr, keys)
	// My Info to Send
	localMinerInfo := MinerInfo{localAddr, m.PublKey}

	//setup an ArtNode Reciever
	artNodeInst := new(ArtNodeInstance)
	// register art node instance locally
	rpc.Register(artNodeInst)

	// Add a listener on myself
	localListener, err := net.ListenTCP("tcp", localAddr)
	checkError(err)

	// Connect to server
	serverConn, err := connectServer(serverAddr, localMinerInfo, m.Settings)
	checkError(err)
	fmt.Println("Settings ", m.Settings)

	// Setup Heartbeats
	go doEvery(time.Duration(m.Settings.HeartBeat-3)*time.Millisecond, serverConn.SendHeartbeat)

	// TODO: Check in with neighbors
	// TODO: Ask Neighbors for blockchain that already exists
	genesisBlock := m.CreateGenesisBlock()

	m.Blockchain = minerlib.NewBlockchainStorage(genesisBlock, m.Settings)
	checkError(err)
	go m.StartMining()
	// go m.TestEarlyExit()

	// Ask for Neighbors
	fmt.Println("Asking for neighbours")
	err = serverConn.RequestMiners(&miners, m.Settings.MinNumMinerConnections)
	fmt.Println("Got neighbours!")
	checkError(err)
	// Set up RPC server for other miners to hitng miners"))
	fmt.Println("miners1: ", miners)

	err = m.AddMinersToList(&miners)
	checkError(err)
	fmt.Printf("miners1: %v \n", &m.Neighbors)

	err = m.OpenNeighborConnections()
	checkError(err)
	fmt.Println("Opened RPC connections to neighbor miners")

	neighborToReceiveBCFrom, err := m.ConnectToNeighborMiners(localAddr)
	checkError(err)
	fmt.Println("Connected to neighbor miners; received BCTree")

	// Set up receiver for other Miners
	minerReceiverInst := new(MinerReceiverInstance)
	rpc.Register(minerReceiverInst)

	serviceRequests(localListener)
}

func connectServer(serverAddr *net.TCPAddr, minerInfo MinerInfo, settings *blockartlib.MinerNetSettings) (serverConnection minerlib.ServerInstance, err error) {
	// dial to server
	serverRPCClient, err := rpc.Dial("tcp", serverAddr.String())
	checkError(err)
	// setup gob
	gob.Register(&net.TCPAddr{})
	gob.Register(&elliptic.CurveParams{})

	//1st rpc call
	//2nd retrieve settings ==> 2 in 1
	err = serverRPCClient.Call("RServer.Register", minerInfo, settings)
	checkError(err)
	// Create the serverConnection.
	// TODO: refactor to ServerInstance
	tcpFromAddr, err := net.ResolveTCPAddr("tcp", minerInfo.Address.String())
	checkError(err)
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
		checkError(err)
		defer conn.Close()

		go rpc.ServeConn(conn)
		fmt.Println("after connection served")

		time.Sleep(10 * time.Millisecond)

		if len(OpQueue) != 0 {
			fmt.Println("connect to queue")
			artNodeConnector, err = rpc.Dial("tcp", OpQueue[0].LocalIP)
			checkError(err)
			var b bool
			err = artNodeConnector.Call("MinerInstance.ConnectMiner", true, &b)
			checkError(err)
			OpQueue = OpQueue[1:]
			fmt.Println("connected to queue ", b, "len ", len(OpQueue))
		}
		// DrawOperations to validate
		// For valid add to miner op channel
	}
}

func checkError(err error) {
	if err != nil {
		// induce a panic for stacktrace
		fmt.Println("There was an error: ")
		fmt.Println(err)
		/*a := []string{"hola"}
		fmt.Printf("Error: %v %v\n", err, a[32])*/
		os.Exit(1)
	}
}

func doEvery(d time.Duration, f func(time.Time) error) error {
	for x := range time.Tick(d) {
		f(x)
	}
	return nil
}

func getKeyPair(pubStr string, privStr string) (*blockartlib.KeyPair, error) {
	// TODO: Fix w/e is up with unicode vs strings
	//priv, pub := keys.Decode(privStr, pubStr)
	priv, pub, err := keys.Generate()
	checkError(err)
	pair := blockartlib.KeyPair{
		Private: priv,
		Public:  pub,
	}
	return &pair, nil
}

// =========================
// Connection Instances
// TODO: Extract out from ink-miner.go
// =========================

type ArtNodeInstance int // same as above

func (si *ArtNodeInstance) ConnectNode(an *blockartlib.ArtNodeInstruction, reply *bool) error {
	fmt.Println("In rpc call to register the AN")
	privateKey := keys.DecodePrivateKey(an.PrivKey)
	publicKey := keys.DecodePublicKey(an.PubKey)
	if !keys.MatchPrivateKeys(privateKey, m.PrivKey) && !keys.MatchPublicKeys(publicKey, m.PublKey) {

		fmt.Println("Private keys do not match.")
		return blockartlib.DisconnectedError("Key pair isn't valid")
	} else {
		*reply = true
		OpQueue = append(OpQueue, an)
	}
	return nil
}

func (si *ArtNodeInstance) GetGenesisBlockHash(stub *bool, reply *string) error {
	fmt.Println("In RPC getting hash of genesis block")
	*reply = m.Settings.GenesisBlockHash
	return nil
}

func (si *ArtNodeInstance) GetAvailableInk(stub *bool, reply *uint32) error {
	fmt.Println("In RPC getting ink from miner")
	*reply = m.InkLevel
	return nil
}

func (si *ArtNodeInstance) GetBlockChildren(hash *string, reply *[]string) error {
	fmt.Println("In RPC getting children hashes")
	// bla, err := m.Blockchain.GetChildrenNodes(*hash)
	// *reply = bla
	// return err
	return nil
}

type MinerReceiverInstance int

func (si *MinerReceiverInstance) ConnectNewNeighbor(neighborAddr *net.TCPAddr, reply *int) error {
	// Add neighbor to list of neighbors
	var newNeighbor = minerlib.MinerConnection{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", neighborAddr.String())
	if err != nil {
		return err
	}
	newNeighbor.Addr = *tcpAddr
	m.Neighbors = append(m.Neighbors, &newNeighbor)

	// Return the length of our blockchain, so the new neighbor can decide
	// if they want our tree.
	*reply = m.Blockchain.BC.LastNode.Current.Depth

	return nil
}

// struct for communicating info about a miner to the server
type MinerInfo struct {
	Address net.Addr
	Key     *ecdsa.PublicKey
}
