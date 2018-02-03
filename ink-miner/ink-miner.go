package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strings"
	"minerlib"
	"net/rpc"
	"blockartlib"
	"net"
	"time"
)

var m minerlib.Miner

type ServerInstance int // for now it's the int, but we can change to actual struct
type ArtNodeInstance int // same as above

var serverConnector *rpc.Client
var artNodeConnector *rpc.Client

func (si *ArtNodeInstance) ConnectNode ( an *blockartlib.ArtNode , reply *bool) error {
	// TODO: miner must be a global variable, so we could call all methods on this 1 instance
	m.ValidateNewArtIdent() // TODO: implement correct arguments, function itself, and return values
	return nil
}

func main() {
	keys := GetKeyPair() // Temporary, keys will be passed in as command line args
	// Need to print then pass to client
	// Connect to server
	localIP := "127.0.0.1:0"
	// TODO: change TBD when it will be available
	serverAddr, err := net.ResolveTCPAddr("tcp", "")
	//m := minerlib.Miner { <-- was local became global
	config := blockartlib.MinerNetSettings{}
	m, err := minerlib.NewMiner(serverAddr, keys, &config)


	fmt.Printf("miner ip: %v, m: %v, \n", localIP, m)

	//setup an RPC connection with AN

	serverInst := new(ServerInstance)
	artNodeInst := new(ArtNodeInstance)

	rpc.Register(serverInst)
	rpc.Register(artNodeInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", localIP)
	CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)

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

func GetKeyPair() *blockartlib.KeyPair {
	curve := elliptic.P256()
	r := strings.NewReader("Hello, Reader!")
	keys, _ := ecdsa.GenerateKey(curve, r)
	fmt.Printf("Keys: %v\n", keys.PublicKey)
	return nil
}
