package minerlib

import (
	"blockartlib"
	"fmt"
	"crypto/ecdsa"
	"net"
	"keys"
)

type Blockchain struct {
	GenesisNode *BCTreeNode
	Head *BCTreeNode
}

// signal channels
var finished chan struct{}

// pipeline channels
var emptyBlocks chan *Block
var operationQueue chan *blockartlib.Operation

const (
	OP_THRESHOLD = 4
	MAX_WAITING_OPS = 10
	MAX_EMPTY_BLOCKS = 3
)

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

// Miner constructor
func NewMiner(serverAddr *net.TCPAddr, keys *blockartlib.KeyPair) (miner Miner) {
	var m = Miner{
		InkLevel: 0,
		ServerNodeAddr: serverAddr,
		ServerHrtBtAddr: nil,
		ArtNodes: []*ArtNodeConnection{},
		Neighbors: []*MinerConnection{},
		PublKey: keys.Public,
		PrivKey: keys.Private,
		Chain: Blockchain{},
		Settings: &blockartlib.MinerNetSettings{},
		LocalCanvas: CanvasData{},
	}
	return m
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

func (m *Miner) StartMining() (err error) {
	// setup channels
	operationQueue = make(chan *blockartlib.Operation, MAX_WAITING_OPS)
	finished  = make(chan struct{})
	emptyBlocks = make(chan *Block, MAX_EMPTY_BLOCKS)
	go m.CreateEmptyBlocks()
	go m.MineBlocks()
	return nil
}

func (m *Miner) CreateEmptyBlocks() {
	for {
		select {
		case <- finished:
			// Done
			return
		default:
			newBlock := &Block{
				ParentHash: m.Chain.Head.BlockHash,
				MinerPublicKey: m.PublKey,
			}
			emptyBlocks <- newBlock
		}
	}
}

func (m *Miner) MineBlocks() (err error) {
	for {
		select {
		case <- finished:
			// done
			return nil
		default:
			b := <- emptyBlocks
			difficulty := uint8(0)
			if len(operationQueue) >= OP_THRESHOLD {
				difficulty = m.Settings.PoWDifficultyOpBlock
				for i := 0; i <= OP_THRESHOLD; i++ {
					b.Operations = append(b.Operations, <- operationQueue)
				}
			} else {
				difficulty = m.Settings.PoWDifficultyNoOpBlock
			}
			err = b.Mine(difficulty)
			if err != nil {
				return err
			}
		}
	}
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

