package minerlib

import (
	"blockartlib"
	"time"
	"net"
	"net/rpc"
)

type Blockchain struct {
	GenesisNode *BCTreeNode
	// perhaps longest chain, addable block, etc.
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

type Miner struct {
	InkLevel int
	ServerNodeAddr *net.TCPAddr
	ServerHrtBtAddr *net.TCPAddr
	ArtNodes []*ArtNodeConnection
	Neighbors []*MinerConnection
	PublKey string
	PrivKey string
	Chain Blockchain
}


/*
func (m *Miner) callAll() {
	for _, artNode := range m.ArtNodes {
		artNode.RPCClient, err = rpc.Dial(artNode.Addr)
		err = artNode.RPCClient.Call("All", arg, response)
	}
}
*/


func NewMiner(serverAddr net.TCPAddr, keys *blockartlib.KeyPair, config *blockartlib.MinerNetSettings) (miner *Miner, err error) {
	return nil, nil
}

func (m *Miner) ValidateNewArtIdent() (err error) {
	return nil
}

func (m *Miner) IsEnoughInk() (err error) {
	return nil
}

func (m *Miner) GenerateNoopBlock() (err error) {
	return nil
}

func (m *Miner) GenerateOpBlock() (err error) {
	return nil
}

// validates incoming block from other miner
func (m *Miner) ValidateBlock() (err error){
	// TODO: include here check against the block produced (or paused?)
	// if block arrived during generating process
	// or before sending the generated block out ===> TODO: DOUBLE SPENDING CHECK
	return nil
}

// this should pause (or delete?) process of mining ink
// in case we want to keep the point at which we've stopped
func (m *Miner) StopMining() (err error){
	return nil
}

// this should resume process of mining ink
func (m *Miner) ResumeMining() (err error){
	return nil
}

/////// functions to perform operations on the blockchain

func (m *Miner) AddBlockToBC() (err error){
	return nil
}

func (m *Miner) RemoveBlockFromBC() (err error){
	return nil
}

// func (m *Miner) FetchParent() (b *Block

/////// functions to interact with server

// retrieves settings from server
func (m *Miner) RetrieveSettings() (err error) {
	return nil
}

// requests another miner's ID (info) from the server
func (m *Miner) RequestMiner() (err error) {
	return nil
}

func (m *Miner) SendHeartbeat(t time.Time) (err error) {
	//var result bool
	// TODO: make it correct once server will be alive
/*	erro := rpc.Client{}.Call("RPC call on server", m.PubKey, &result)
	if erro !=nil || !result{
		return blockartlib.DisconnectedError("error")
	}*/
	// TODO: stop here
	return nil
}

/////// functions to interact with other miners

func (m *Miner) ConnectToOtherMiner() (err error) {
	return nil
}

// sends out the block to other miners
func (m *Miner) DisseminateBlock() (err error) {
	return nil
}

/////// helpers

func (m *Miner) AddInk() (err error) {
	return nil
}

func (m *Miner) DrawInk() (err error) {
	return nil
}
