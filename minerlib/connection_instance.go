package minerlib

import (
	"net"
	"net/rpc"
	"time"
	"crypto/ecdsa"
	"encoding/gob"
	"crypto/elliptic"
)

type ServerInstance struct {
	Addr net.TCPAddr
	RPCClient *rpc.Client
	Public *ecdsa.PublicKey
}

// requests another miner's ID (info) from the server
func (m *ServerInstance) RequestMiners(lom *[]net.Addr, minNeighbours uint8) (err error) {
	gob.Register(&net.TCPAddr{})
	gob.Register(&[]net.Addr{})
	gob.Register(&[]net.TCPAddr{})
	gob.Register(&elliptic.CurveParams{})

	//for uint8(len(*lom))<minNeighbours { // TODO: uncomment for production
	for uint8(len(*lom))<2 {
		//fmt.Println("request lom")
		err = m.RPCClient.Call("RServer.GetNodes", m.Public, &lom)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ServerInstance) SendHeartbeat(t time.Time) (err error) {
	var ignored bool
	err = m.RPCClient.Call( "RServer.HeartBeat", m.Public, &ignored)
	if err != nil {
		return err
	}
	return nil
}

// Structs to manage the connections to other entities
type ArtNodeConnection struct {
	// Addr to Dial
	Addr net.TCPAddr
	// RPC Client to Call
	RPCClient *rpc.Client
}

type MinerConnection struct {
	Addr net.TCPAddr
	RPCClient *rpc.Client
}
