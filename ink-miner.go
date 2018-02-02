package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strings"

	"./minerlib"
	"net/rpc"
	"project1/blockartlib"
	"net"
	"time"
)

type ServerInst int // for now it's the int, but we can change to actual struct
type ArtNodeInst int // same as above

type ArtNode struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
	LocalIP		string
}

var m minerlib.Miner // TODO: <--- introduced globally instead local var in main()

var serverConnector *rpc.Client
var artNodeConnector *rpc.Client

func (si *ArtNodeInst) ConnectNode ( an *ArtNode , reply *bool) error {
	// TODO: miner must be a global variable, so we could call all methods on this 1 instance
	new(minerlib.Miner).ValidateNewArtlib() // TODO: implement correct arguments, function itself, and return values
	return nil
}

func main() {
	GetKeyPair() // Temporary, keys will be passed in as command line args
	// Need to print then pass to client
	// Connect to server
	localIP := "127.0.0.1:0"
	//serverAddr := "tbd"  // TODO: change TBD when it will be available
	var nbrs [256]int
	//m := minerlib.Miner { <-- was local became global
	m = minerlib.Miner {
		nbrs,
		"server addr",
		false,
		nil,
		"pubKey",
		"privKey",
		make([]int, 1),
		new(&blockartlib.MinerNetSettings{}),
		}

	fmt.Printf("miner ip: %v, m: %v, \n", localIP, m)

	//setup an RPC connection with AN

	serverInst := new(ServerInst)
	artNodeInst := new(ArtNodeInst)

	rpc.Register(serverInst)
	rpc.Register(artNodeInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", localIP)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		CheckError(err)
		defer conn.Close()

		go rpc.ServeConn(conn)

		// sending heartbeat for every X seconds
		var x int = int(m.Settings.HeartBeat)
		go doEvery(time.Duration(x) * time.Second, m.SendHeartbeat)

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

func GetKeyPair() {
	curve := elliptic.P256()
	r := strings.NewReader("Hello, Reader!")
	keys, _ := ecdsa.GenerateKey(curve, r)
	fmt.Printf("Keys: %v\n", keys.PublicKey)
}
