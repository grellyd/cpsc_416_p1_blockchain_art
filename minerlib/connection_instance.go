package minerlib

import (
	"net"
	"net/rpc"
	"crypto/ecdsa"
)

type MinerCaller struct {
	Addr net.TCPAddr
	RPCClient *rpc.Client
	Public *ecdsa.PublicKey
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

