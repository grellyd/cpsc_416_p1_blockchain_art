package minerlib

import (
	"blockartlib"
	"crypto/ecdsa"
	"fmt"
	"minerlib/compute"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// signal channels
var doneMining chan struct{}
var earlyExitSignal chan struct{}
var exited chan struct{}

// waitgroups
var minersGroup sync.WaitGroup

const (
	OP_THRESHOLD     = 4
	MAX_WAITING_OPS  = 10
	MAX_EMPTY_BLOCKS = 3
	NUM_MINING_TASKS = 1
)

// maps hashes to blocks for the invalid blocks
type Forest map[string]*Block

type Miner struct {
	InkLevel        uint32
	ServerNodeAddr  *net.TCPAddr
	ServerHrtBtAddr *net.TCPAddr
	ArtNodes        []*ArtNodeConnection
	Neighbors       []*MinerConnection
	PublKey         *ecdsa.PublicKey
	PrivKey         *ecdsa.PrivateKey
	Blockchain      *BCStorage
	Settings        *blockartlib.MinerNetSettings
	LocalCanvas     CanvasData
	BlockForest     map[string]*Block
	// pipeline channel
	operationQueue chan *blockartlib.Operation
}

// Miner constructor
func NewMiner(serverAddr *net.TCPAddr, keys *blockartlib.KeyPair) (miner Miner) {
	var m = Miner{
		InkLevel:        0,
		ServerNodeAddr:  nil,
		ServerHrtBtAddr: serverAddr,
		ArtNodes:        []*ArtNodeConnection{},
		Neighbors:       []*MinerConnection{},
		PublKey:         keys.Public,
		PrivKey:         keys.Private,
		Blockchain:      nil,
		Settings:        &blockartlib.MinerNetSettings{},
		LocalCanvas:     CanvasData{},
		BlockForest:     map[string]*Block{},
		operationQueue: make(chan *blockartlib.Operation, MAX_WAITING_OPS),
	}
	return m
}

func (m *Miner) IsEnoughInk() (err error) {
	return nil
}

func (m *Miner) AddOp(o *blockartlib.Operation) error {
	// blocks until space in buffered channel
	m.operationQueue <- o
	return nil
}

func (m *Miner) CreateGenesisBlock() (g *Block) {
	return NewBlock(m.Settings.GenesisBlockHash, nil)
}

func (m *Miner) StartMining() (err error) {
	fmt.Printf("[miner] Starting Mining Process\n")
	// setup channels
	doneMining = make(chan struct{})
	earlyExitSignal = make(chan struct{})
	for i := 0; i < NUM_MINING_TASKS; i++ {
		go m.MineBlocks()
		minersGroup.Add(1)
	}
	return nil
}

func (m *Miner) TestEarlyExit() {
	time.Sleep(10000 * time.Millisecond)
	fmt.Printf("[miner] Killing...\n")
	err := m.StopMining()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (m *Miner) MineBlocks() (err error) {
	fmt.Printf("[miner] Starting to Mine Blocks\n")
	for {
		select {
		case <-doneMining:
			// done
			fmt.Printf("[miner] Done Mining Blocks\n")
			return nil
		case <-earlyExitSignal:
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
			if len(m.operationQueue) >= OP_THRESHOLD {
				difficulty = m.Settings.PoWDifficultyOpBlock
				for i := 0; i <= OP_THRESHOLD; i++ {
					// TODO: Check Valid operations
					b.Operations = append(b.Operations, <- m.operationQueue)
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
func (m *Miner) ValidBlock(b *Block) (valid bool, err error) {
	hash, err := b.Hash()
	if err != nil {
		return false, fmt.Errorf("Unable validate block: %v", err)
	}
	difficulty := uint8(0)
	if len(b.Operations) > 0 {
		difficulty = m.Settings.PoWDifficultyOpBlock
		// check each op has a valid sig
		for _, op := range b.Operations {
			if op.ShapeHash != op.CalculateSig() {
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
func (m *Miner) StopMining() (err error) {
	fmt.Printf("[miner] Stopping Mining by command\n")
	close(earlyExitSignal)
	minersGroup.Wait()
	fmt.Printf("[miner] Stopped\n")
	return nil
}

/////// functions to interact with other miners

func (m *Miner) OpenNeighborConnections() (err error) {
	/* Opens RPC connections to neighbouring Miners and fills in the
	   RPCClient field in the corresponding MinerConnection struct */
	for i, neighbor := range m.Neighbors {
		if neighbor.Alive {
			continue
		}
		fmt.Println("Before open RPC to neighbour: ", neighbor.Addr.String())
		neighbor.RPCClient, err = rpc.Dial("tcp", neighbor.Addr.String())
		if err != nil {
			deleteNeighbour(m, i)
			return nil
		}
		fmt.Printf("Opened RPC connection to neighbor with tcpAddr %s\n", neighbor.Addr.String())
	}

	return nil
}

// TODO: Actually handle the case where the blockchain we choose is invalid
func (m *Miner) ConnectToNeighborMiners(localAddr *net.TCPAddr) (bestNeighbor net.TCPAddr, err error) {
	/* Makes the RPC call to register itself to neighbouring miners.
	   Neighbours will respond with the length of their Blockchain;
	   does NOT currently account for the fact that the given chain
	   may be invalid.

	   RETURNS: the net.TCPAddress of the neighbour with the longest
	   blockchain depth
	*/

	// Connect to each neighbour miner and keep track of the one with the largest depth
	var bestMinerAddr net.TCPAddr
	largestDepth := 0
	depth := 0

	fmt.Println("Our address before sending RPC call: ", localAddr.String())
	for i, connection := range m.Neighbors {
		/*fmt.Println("DISCONNECT!!!")
		time.Sleep(4*time.Second)*/
		err = connection.RPCClient.Call("MinerInstance.ConnectNewNeighbor", localAddr, &depth)
		if err != nil {
			// TODO: Should we just ignore this miner then and move on to the next one?
			erro := deleteNeighbour(m, i)
			blockartlib.CheckErr(erro)
			//return net.TCPAddr{}, err
			continue
		}

		if (depth >= largestDepth) {
			largestDepth = depth
			bestMinerAddr = connection.Addr
		}
	}

	return bestMinerAddr, nil
}

// TODO: Actually handle the case where the blockchain we choose is invalid
func (m *Miner) RequestBCStorageFromNeighbor(neighborAddr *net.TCPAddr) (err error) {
	treeArray := make([][]byte, 0)
	for i, v := range m.Neighbors {
		if v.Addr.String() == neighborAddr.String() {
			err := v.RPCClient.Call("MinerInstance.DisseminateTree", true, &treeArray)
			if err != nil {
				deleteNeighbour(m, i)
			}
		}
	}
	if len(treeArray) != 0 {
		reconstructTree(m, &treeArray)
	}

	return nil
}

// sends out the block to other miners
func (m *Miner) DisseminateBlock(block *Block) (err error) {
	// TODO: notify waiting artNodes if your block is op number of nodes deep now
	// TODO: Not sure this is the right spot for this.
	for _,v := range m.Neighbors {
		marshalledBlock, err := block.MarshallBinary()
		blockartlib.CheckErr(err)
		var b bool
		err = v.RPCClient.Call("MinerInstance.ReceiveBlockFromNeighbour", &marshalledBlock, &b)
		if !b {
			fmt.Println("Bad block") // TODO: think what to do in this case
		}
	}
	return err
}

func (m *Miner) DisseminateOperation(op Operation) (err error) {
	return err
}

/////// helpers

func (m *Miner) AddInk() (err error) {
	return nil
}

func (m *Miner) DrawInk() (err error) {
	return nil
}

func (m *Miner) IsMinerInList() (err error) {
	return nil
}

func (m *Miner) MarshallTree (result *[][]byte) {
	tree := m.Blockchain.BCT
	//genBlock := tree.GenesisNode.BlockResiding
	//marshalledGenBlock, _ := genBlock.MarshallBinary()
	//*result = append(*result, marshalledGenBlock)
	//for len(tree.GenesisNode.Children) != 0 {
	for range tree.GenesisNode.Children {
		children, err := m.Blockchain.GetChildrenNodes(tree.GenesisNode.CurrentHash)
		blockartlib.CheckErr(err)
		for _, v := range children {
			node := FindBCTreeNode(tree.GenesisNode,v)
			block := node.BlockResiding
			marshalledBlock, err := block.MarshallBinary()
			if err != nil {
				fmt.Println("error happened: ", err)
				continue
			}
			*result = append(*result, marshalledBlock)
		}
	}
	return
}

func (m *Miner) AddMinersToList(lom *[]net.Addr) (err error) {
	if len(*lom) == 0 {
		return nil
	} else if len(m.Neighbors) == 0 {
		for _, val := range *lom {
			addMinerToList(m, val)
		}
	} else if len(m.Neighbors) > 0 {
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

func addMinerToList(m *Miner, addr net.Addr) error {
	var newNeighbor = MinerConnection{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return err
	}
	newNeighbor.Addr = *tcpAddr
	m.Neighbors = append(m.Neighbors, &newNeighbor)
	return nil
}

func isMinerInList(m *Miner, addr net.Addr) bool {
	for _, v := range m.Neighbors {
		if v.Addr.String() == addr.String() {
			return true
		}
	}
	return false
}

func deleteNeighbour (m *Miner, index int) error {
	buf := m.Neighbors[:index]
	m.Neighbors = append(buf, m.Neighbors[index+1:]...)
	return nil
}

func reconstructTree(m *Miner, tree *[][]byte) {
	t := *tree
	genBlock := m.CreateGenesisBlock()
	fmt.Println("tree: ", t)
	m.Blockchain = NewBlockchainStorage(genBlock, m.Settings)
	t = t[1:]
	for _,v := range t {
		b, err := UnmarshallBinary(v)
		if err != nil {
			fmt.Println("unmarshalling failed")
			return
		}
		m.Blockchain.AppendBlock(b, m.Settings)
	}

}
