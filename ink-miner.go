package main

import (

	"fmt"
	"minerlib"
	"net/rpc"
	"blockartlib"
	"net"
	"time"
	"os"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
)
var m minerlib.Miner // singleton for miner
var bc minerlib.BCStorage //singleton for blockchain stored at miner

type ServerInstance int // for now it's the int, but we can change to actual struct
type ArtNodeInstance int // same as above

var serverConnector *rpc.Client
var artNodeConnector *rpc.Client
var artNodeQueue []*blockartlib.ArtNodeInstruction

func (si *ArtNodeInstance) ConnectNode ( an *blockartlib.ArtNodeInstruction , reply *bool) error {
	fmt.Println("In rpc call to register the AN")
	err := m.ValidateNewArtIdent(an)
	CheckError(err)
	if err == nil {
		*reply = true
		artNodeQueue = append(artNodeQueue, an)
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

func (si *ArtNodeInstance) GetBlockChildren (hash *string, reply *[]string) error {
	fmt.Println("In RPC getting children hashes")
	bla, err := bc.GetChildrenNodes(*hash)
	*reply = bla
	return err
}

func main() {

	fmt.Println("start")
	args := os.Args[1:]

	// Missing command line args.
	if len(args) != 3 {
		fmt.Println("usage: go run ink-miner.go [server ip:port] [pubKey] [privKey] ")
		return
	}

	servAddr := args[0]
	keys := GetKeyPair(args)
	// Connect to server
	localIP := "127.0.0.1:0"
	// TODO: change TBD when it will be available
	serverAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	config := blockartlib.MinerNetSettings{}
	fmt.Println("config ", config)
	m, err = minerlib.NewMiner(serverAddr, keys, &config)

	fmt.Printf("miner ip: %v, m: %v, \n", localIP, m)

	//setup an RPC connection with AN

	//serverInst := new(ServerInstance)
	artNodeInst := new(ArtNodeInstance)

	//rpc.Register(serverInst)
	rpc.Register(artNodeInst)

	tcpAddr, err := net.ResolveTCPAddr("tcp", localIP)
	CheckError(err)


	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)
	fmt.Println("TCP address: ", listener.Addr().String())

	serverConnector, err = rpc.Dial("tcp", servAddr)

	err = m.ConnectServer(serverConnector, listener.Addr().String())
	CheckError(err)

	var mConn minerlib.MinerCaller
	mConn.Addr = *tcpAddr
	mConn.RPCClient = serverConnector
	mConn.Public = m.PublKey

	go doEvery(time.Duration(m.Settings.HeartBeat-3) * time.Millisecond, mConn.SendHeartbeat)

	var lom = []net.Addr{}

	// TODO: put it into a goroutine wh will as the miners until it will fill the entire array
	err = mConn.RequestMiner(&lom, m.Settings.MinNumMinerConnections)
	fmt.Println("LOM1: ", lom)
	CheckError(err)

	err = m.AddMinersToList(&lom)
	fmt.Printf("LOM1: %v \n", &m.Neighbors)
	CheckError(err)

	for {
		conn, err := listener.Accept()
		CheckError(err)
		defer conn.Close()

		go rpc.ServeConn(conn)
		fmt.Println("after connection served")

		time.Sleep(10*time.Millisecond)

		if len(artNodeQueue) != 0{
			fmt.Println("connect to queue")
			artNodeConnector, err = rpc.Dial("tcp", artNodeQueue[0].LocalIP)
			CheckError(err)
			var b bool
			err = artNodeConnector.Call("MinerInstance.ConnectMiner", true, &b)
			CheckError(err)
			artNodeQueue = artNodeQueue[1:]
			fmt.Println("connected to queue ", b, "len ", len(artNodeQueue))
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

func GetKeyPair(args[] string) *blockartlib.KeyPair {
	if len(args) == 0 {
		return nil
	}
	a := args[2]
	b := args[1]
	priv, pub := decodeKeys(a, b)
	pair := blockartlib.KeyPair{
		priv,
		pub,
	}
	return &pair
}

func decodeKeys(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {

	c, _ := hex.DecodeString(pemEncoded)
	d, _ := hex.DecodeString(pemEncodedPub)
	privateKey, _ := x509.ParseECPrivateKey(c)
	fmt.Println(privateKey)
	genericPublicKey, _ := x509.ParsePKIXPublicKey(d)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	//publicKey := &privateKey.PublicKey

	return privateKey, publicKey
}


