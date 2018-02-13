package main

import (
	"os"
	"fmt"
	"minerlib"
	"net/rpc"
	"blockartlib"
	"net"
	"time"
	"keys"
	"crypto/ecdsa"
	"encoding/gob"
	"crypto/elliptic"
)
var m minerlib.Miner
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
	CheckError(err)
	serverAddr, err := net.ResolveTCPAddr("tcp", args[0])
	CheckError(err)
	localAddr, err := net.ResolveTCPAddr("tcp", localIP)
	CheckError(err)
	
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
	CheckError(err)
	
	// Connect to server
	// TODO: Check settings are working as expected
	serverConn, err := connectServer(serverAddr, localMinerInfo, m.Settings)
	CheckError(err)
	
	// Setup Heartbeats
	go doEvery(time.Duration(m.Settings.HeartBeat-3) * time.Millisecond, serverConn.SendHeartbeat)

	// TODO:
	// go Mine NOPS

	// Ask for Neighbors
	err = serverConn.RequestMiners(&miners, m.Settings.MinNumMinerConnections)
	CheckError(err)
	fmt.Println("miners1: ", miners)

	err = m.AddMinersToList(&miners)
	CheckError(err)
	fmt.Printf("miners1: %v \n", &m.Neighbors)
	serviceRequests(localListener)
}


func connectServer(serverAddr *net.TCPAddr, minerInfo MinerInfo, settings *blockartlib.MinerNetSettings) (serverConnection minerlib.MinerCaller, err error) {
	// dial to server
	serverRPCClient, err := rpc.Dial("tcp", serverAddr.String())
	// setup gob
	gob.Register(&net.TCPAddr{})
	gob.Register(&elliptic.CurveParams{})

	//1st rpc call
	//2nd retrieve settings ==> 2 in 1
	err = serverRPCClient.Call("RServer.Register", minerInfo, settings)
	CheckError(err)
	fmt.Println("Settings ", settings)
	// Create the serverConnection. 
	// TODO: refactor to ServerInstance
	serverConnection = minerlib.MinerCaller{
		Addr: *minerInfo.Address,
		RPCClient: serverConnector,
		Public: minerInfo.Key,
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

		time.Sleep(10*time.Millisecond)

		if len(OpQueue) != 0{
			fmt.Println("connect to queue")
			artNodeConnector, err = rpc.Dial("tcp", OpQueue[0].LocalIP)
			CheckError(err)
			var b bool
			err = artNodeConnector.Call("MinerInstance.ConnectMiner", true, &b)
			CheckError(err)
			OpQueue = OpQueue[1:]
			fmt.Println("connected to queue ", b, "len ", len(OpQueue))
		}
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	os.Exit(1)
}

func doEvery(d time.Duration, f func(time.Time) error) error {
	for x := range time.Tick(d) {
		f(x)
	}
	return nil
}

func getKeyPair(pubStr string, privStr string) (*blockartlib.KeyPair, error) {
	priv, pub := keys.Decode(privStr, pubStr)
	pair := blockartlib.KeyPair{
		Private: priv,
		Public: pub,
	}
	return &pair, nil
}

// =========================
// Connection Instances
// TODO: Extract out from ink-miner.go
// =========================

type ServerInstance int // for now it's the int, but we can change to actual struct
type ArtNodeInstance int // same as above

func (si *ArtNodeInstance) ConnectNode(an *blockartlib.ArtNodeInstruction , reply *bool) error {
	fmt.Println("In rpc call to register the AN")
	err := m.ValidateNewArtIdent(an)
	if err == nil {
		*reply = true
		OpQueue = append(OpQueue, an)
	} else {
		err = fmt.Errorf("Unable to connect to node: %v", err)
	}
	return err
}

func (si *ArtNodeInstance) GetGenesisBlockHash (stub *bool, reply *string) error {
	fmt.Println("In RPC getting hash of genesis block")
	*reply = m.Settings.GenesisBlockHash
	return nil
}

func (si *ArtNodeInstance) GetAvailableInk (stub *bool, reply *uint32) error {
	fmt.Println("In RPC getting ink from miner")
	*reply = m.InkLevel
	return nil
}

// struct for communicating info about a miner to the server
type MinerInfo struct {
	Address *net.TCPAddr
	Key     *ecdsa.PublicKey
}
