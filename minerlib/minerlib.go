package minerlib

import (
	"blockartlib"
	"time"
	"fmt"
	"crypto/ecdsa"
	"encoding/gob"
	"crypto/elliptic"
	"net"
	"net/rpc"
	"keys"
)

type Blockchain struct {
	GenesisNode *BCTreeNode
	// perhaps longest chain, addable block, etc.
}

type Miner struct {
	InkLevel uint32
	ServerNodeAddr *net.TCPAddr
	ServerHrtBtAddr *net.TCPAddr
	ArtNodes []*ArtNodeConnection
	Neighbors []*MinerConnection
	PublKey *ecdsa.PublicKey
	PrivKey *ecdsa.PrivateKey
	Chain Blockchain
	Settings *blockartlib.MinerNetSettings
	LocalCanvas CanvasData
}

type MinerInfo struct {
	Address net.Addr
	Key     ecdsa.PublicKey
}

// Miner constructor
func NewMiner(serverAddr *net.TCPAddr, keys *blockartlib.KeyPair, config *blockartlib.MinerNetSettings) (miner Miner, err error) {
	var m = Miner{
		InkLevel: 0,
		ServerNodeAddr: serverAddr,
		ServerHrtBtAddr: nil,
		ArtNodes: []*ArtNodeConnection{},
		Neighbors: []*MinerConnection{},
		PublKey: keys.Public,
		PrivKey: keys.Private,
		Chain: Blockchain{},
		Settings: config,
		LocalCanvas: CanvasData{},
	}
	return m, nil
}

func (m *Miner) ValidateNewArtIdent(an *blockartlib.ArtNodeInstruction) (err error) {
	privateKey := keys.DecodePrivateKey(an.PrivKey)
	publicKey := keys.DecodePublicKey(an.PubKey)
	if !keys.MatchPrivateKeys(privateKey, m.PrivKey) && !keys.MatchPublicKeys(publicKey, m.PublKey) {

		fmt.Println("Private keys do not match.")
		return blockartlib.DisconnectedError("Key pair isn't valid")
	}
	fmt.Println("keys match")
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

func (m *Miner) MineBlock(b *Block) (err error) {
	difficulty := uint8(0)
	if len(b.Operations) == 0 {
		difficulty = m.Settings.PoWDifficultyNoOpBlock
	} else {
		difficulty = m.Settings.PoWDifficultyOpBlock
	}
	return b.Mine(difficulty)
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
func (m *Miner) ConnectServer(servConnector *rpc.Client, minerAddr string) (err error) {
	//1st rpc call
	//2nd retrieve settings ==> 2 in 1
	gob.Register(&net.TCPAddr{})
	gob.Register(&elliptic.CurveParams{})
	puk := *m.PublKey
	a, _ := net.ResolveTCPAddr("tcp",minerAddr)
	var minerInfo = MinerInfo{
		net.Addr(a),
		puk,
	}

	err = servConnector.Call("RServer.Register", minerInfo, m.Settings)
	blockartlib.CheckErr(err)
	fmt.Println("Settings ", m.Settings)
	return nil
}

func (m *Miner) RetrieveSettings() (err error) {
	return nil
}

// requests another miner's ID (info) from the server
func (m *MinerCaller) RequestMiner(lom *[]net.Addr, minNeighbours uint8) (err error) {
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

func (m *MinerCaller) SendHeartbeat(t time.Time) (err error) {
	var ignored bool
	err = m.RPCClient.Call( "RServer.HeartBeat", m.Public, &ignored)
	if err != nil {
		return err
	}
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

func (m *Miner) IsMinerInList () (err error) {
	return nil
}

func (m *Miner) AddMinersToList (lom *[]net.Addr) (err error) {
	if len(*lom) == 0 {
		return nil
	} else if len(m.Neighbors) == 0 {
		for _, val := range *lom {
			addMinerToList(m, val)
		}
	}else if len(m.Neighbors) > 0 {
		for _, val := range *lom {
			if len(m.Neighbors) == 256 {
				return nil
			}
			if !isMinerInList(m, val) {
				addMinerToList(m, val)
			}
		}
	}
	return nil
}

func addMinerToList (m *Miner, addr net.Addr) error {
	var newNeighbour = MinerConnection{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return err
	}
	newNeighbour.Addr = *tcpAddr
	m.Neighbors = append(m.Neighbors, &newNeighbour)
	return nil
}

func isMinerInList (m *Miner, addr net.Addr) bool {
	for _, v := range m.Neighbors {
		if v.Addr.String() == addr.String() {
			return true
		}
	}
	return false
}

