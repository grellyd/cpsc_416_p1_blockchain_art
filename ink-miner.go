package main

import (
	"fmt"
	"minerlib"
	"net/rpc"
	"blockartlib"
	"net"
	"time"
	"os"
	"keys"
)
var m minerlib.Miner

var serverConnector *rpc.Client
var artNodeConnector *rpc.Client
var OpQueue []*blockartlib.ArtNodeInstruction

func main() {
	fmt.Println("start")
	args := os.Args[1:]

	// Missing command line args.
	if len(args) != 3 {
		fmt.Println("usage: go run ink-miner.go [server ip:port] [pubKey] [privKey] ")
		return
	}
	servAddr := args[0]
	keys := getKeyPair(args[2], args[1])
	serverListener, serverConn := setupServerConn(servAddr, keys)

	go doEvery(time.Duration(m.Settings.HeartBeat-3) * time.Millisecond, serverConn.SendHeartbeat)

	var miners = []net.Addr{}
	
	// TODO:
	// go Mine NOPS

	// Ask for Neighbors
	err := serverConn.RequestMiner(&miners, m.Settings.MinNumMinerConnections)
	fmt.Println("miners1: ", miners)
	CheckError(err)

	err = m.AddMinersToList(&miners)
	CheckError(err)
	fmt.Printf("miners1: %v \n", &m.Neighbors)
	serviceRequests()
}

func setupServerConn(servAddr string, keys *blockartlib.KeyPair) (listener net.TCPListener, serverConn minerlib.MinerCaller) {
	// Connect to server
	localIP := "127.0.0.1:0"
	// TODO: change TBD when it will be available
	serverAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	config := blockartlib.MinerNetSettings{}
	fmt.Println("config ", config)
	m, err = minerlib.NewMiner(serverAddr, keys, &config)

	fmt.Printf("miner ip: %v, m: %v, \n", localIP, m)

}

func setupArtNodeRPC() {
	//setup an RPC connection with AN
	artNodeInst := new(ArtNodeInstance)

	rpc.Register(artNodeInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", localIP)
	CheckError(err)

	listener, err = net.ListenTCP("tcp", tcpAddr)
	CheckError(err)
	fmt.Println("TCP address: ", listener.Addr().String())

	serverConnector, err = rpc.Dial("tcp", servAddr)

	err = m.ConnectServer(serverConnector, listener.Addr().String())
	CheckError(err)

	serverConn.Addr = *tcpAddr
	serverConn.RPCClient = serverConnector
	serverConn.Public = m.PublKey
	return listener, serverConn
}

func serviceRequests(listener net.TCPListener) {
	for {
		conn, err := listener.Accept()
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
}

func doEvery(d time.Duration, f func(time.Time) error) error {
	for x := range time.Tick(d) {
		f(x)
	}
	return nil
}

func getKeyPair(pubStr string, privStr string) *blockartlib.KeyPair {
	priv, pub := keys.Decode(privStr, pubStr)
	pair := blockartlib.KeyPair{
		Private: priv,
		Public: pub,
	}
	return &pair
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
