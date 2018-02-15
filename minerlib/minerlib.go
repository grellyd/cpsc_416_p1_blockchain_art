package minerlib

import (
	"blockartlib"
	"fmt"
	"crypto/ecdsa"
	"net"
	"sync"
	"time"
	"minerlib/compute"
)


// signal channels
var doneMining chan struct{}
var earlyExitSignal chan struct{}
var exited chan struct{}

// pipeline channels
var operationQueue chan *blockartlib.Operation

// waitgroups 
var minersGroup sync.WaitGroup

const (
	OP_THRESHOLD = 4
	MAX_WAITING_OPS = 10
	MAX_EMPTY_BLOCKS = 3
	NUM_MINING_TASKS = 1
)

// maps hashes to blocks for the invalid blocks
type Forest map[string]*Block

type Miner struct {
	InkLevel uint32
	ServerNodeAddr *net.TCPAddr
	ServerHrtBtAddr *net.TCPAddr
	ArtNodes []*ArtNodeConnection
	Neighbors []*MinerConnection
	PublKey *ecdsa.PublicKey
	PrivKey *ecdsa.PrivateKey
	Blockchain *BCStorage
	Settings *blockartlib.MinerNetSettings
	LocalCanvas CanvasData
	BlockForest map[string]*Block
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
		Blockchain: nil,
		Settings: &blockartlib.MinerNetSettings{},
		LocalCanvas: CanvasData{},
		BlockForest: map[string]*Block{},
	}
	return m
}

func (m *Miner) IsEnoughInk() (err error) {
	return nil
}

func (m *Miner) CreateGenesisBlock() (g *Block) {
	return NewBlock(m.Settings.GenesisBlockHash, nil)
}

func (m *Miner) StartMining() (err error) {
	fmt.Printf("[miner] Starting Mining Process\n")
	// setup channels
	operationQueue = make(chan *blockartlib.Operation, MAX_WAITING_OPS)
	doneMining  = make(chan struct{})
	earlyExitSignal = make(chan struct{})
	for i := 0; i < NUM_MINING_TASKS; i++{
		go m.MineBlocks()
		minersGroup.Add(1)
	}
	return nil
}

func (m *Miner)TestEarlyExit() {
	time.Sleep(6000 * time.Millisecond)
	fmt.Printf("[miner] Killing...\n")
	err := m.StopMining()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	err = m.StartMining()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (m *Miner) MineBlocks() (err error) {
	fmt.Printf("[miner] Starting to Mine Blocks\n")
	for {
		select {
		case <- doneMining:
			// done
			fmt.Printf("[miner] Done Mining Blocks\n")
			return nil
		case <- earlyExitSignal:
			minersGroup.Done()
			fmt.Printf("[miner] Early Exit\n")
			return nil
		default:
			parentHash, err := m.Blockchain.LastNodeHash()
			if err != nil {
				return fmt.Errorf("Unable to get parent hash: %v", err)
			}
			b := NewBlock(parentHash, m.PublKey)
			difficulty := uint8(0)
			if len(operationQueue) >= OP_THRESHOLD {
				difficulty = m.Settings.PoWDifficultyOpBlock
				for i := 0; i <= OP_THRESHOLD; i++ {
					b.Operations = append(b.Operations, <- operationQueue)
				}
			} else {
				difficulty = m.Settings.PoWDifficultyNoOpBlock
			}

			fmt.Printf("[miner] Starting Mining: %v\n", b)
			err = b.Mine(difficulty)
			hash, err := b.Hash()
			fmt.Printf("[miner] Done Mining: %v with %s\n", b, hash)

			_ = m.Blockchain.AppendBlock(b, m.Settings)
			if err != nil {
				fmt.Printf("MineBlocks created Error: %v", err)
				return err
			}
		}
	}
}

// validates incoming block from other miner
func (m *Miner) ValidBlock(b *Block) (valid bool, err error){
	hash, err := b.Hash()
	if err != nil {
		return false, fmt.Errorf("Unable validate block: %v", err)
	}
	difficulty := uint8(0)
	if len(b.Operations) > 0 {
		difficulty = m.Settings.PoWDifficultyOpBlock
		// check each op has a valid sig
		for _, op := range b.Operations {
			// TODO
			expectedSig := ""
			err = nil
			if err != nil {
				return false, fmt.Errorf("Unable validate block: %v", err)
			}
			if op.OperationSig != expectedSig {
				return false, nil
			}
		}
	} else {
		difficulty = m.Settings.PoWDifficultyNoOpBlock
	}
	// check nonce adds up
	if !compute.Valid(hash, difficulty) {
		return false, nil
	}

	// check previous block is in tree
	present, err := m.Blockchain.BlockPresent(b)
	if err != nil {
		return false, fmt.Errorf("Unable validate block: %v", err)
	}
	if present {
		return true, nil
	} else {
		// failing that, check its ancestors in the forest, and those pass ValidBlock
		if m.BlockForest[b.ParentHash] != nil {
			parentValid, err := m.ValidBlock(m.BlockForest[b.ParentHash])
			if err != nil {
				return false, fmt.Errorf("Unable validate block: %v", err)
			}
			if parentValid {
				return true, nil
			}
		}
		return false, nil
	}
}

// this stops the process of mining blocks
// Commands the lower level threads to stop.
// Waitgroup finishes when these exit
func (m *Miner) StopMining() (err error){
	fmt.Printf("[miner] Stopping Mining by command\n")
	close(earlyExitSignal)
	minersGroup.Wait()
	fmt.Printf("[miner] Stopped\n")
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

func (m *Miner) ConnectToOtherMiners() (err error) {
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

