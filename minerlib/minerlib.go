package minerlib

import (
	"blockartlib"
	"crypto/ecdsa"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
	"encoding/gob"
	"crypto/elliptic"
)


const (
	OP_THRESHOLD     = 1 // TODO: Change this back when we have flags
	MAX_WAITING_OPS  = 10
	MAX_EMPTY_BLOCKS = 3
	NUM_MINING_TASKS = 1
)

// waitgroups
var minersGroup sync.WaitGroup

// maps hashes to blocks for the invalid blocks
type Forest map[string]*Block

type Miner struct {
	InkLevel        uint32
	ServerNodeAddr  *net.TCPAddr
	ServerHrtBtAddr *net.TCPAddr
	ArtNodes        []*ArtNodeConnection
	Neighbours       []*MinerConnection
	PublKey         *ecdsa.PublicKey
	PrivKey         *ecdsa.PrivateKey
	Blockchain      *BCStorage
	Settings        *blockartlib.MinerNetSettings
	LocalCanvas     CanvasData
	BlockForest     map[string]*Block
	
	// signal channels
	doneMining chan struct{}
	earlyExitSignal chan struct{}
	exited chan struct{}

	OpValidateList	[][]*Operation
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
		Neighbours:       []*MinerConnection{},
		PublKey:         keys.Public,
		PrivKey:         keys.Private,
		Blockchain:      nil,
		Settings:        &blockartlib.MinerNetSettings{},
		LocalCanvas:     CanvasData{},
		BlockForest:     map[string]*Block{},
		OpValidateList: 	 [][]*Operation{},
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
	m.doneMining = make(chan struct{})
	m.earlyExitSignal = make(chan struct{})
	for i := 0; i < NUM_MINING_TASKS; i++ {
		go m.MineBlocks()
		minersGroup.Add(1)
	}
	return nil
}

func (m *Miner) MineBlocks() (err error) {
	fmt.Printf("[miner] Starting to Mine Blocks\n")
	for {
		select {
		case <-m.doneMining:
			// done
			minersGroup.Done()
			fmt.Printf("[miner] Done Mining Blocks\n")
			return nil
		case <-m.earlyExitSignal:
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
			// if there exist enough ops waiting
			fmt.Printf("[miner] OperationQueue: %v\n", m.operationQueue)
			fmt.Printf("[miner] len(OperationQueue): %v\n", len(m.operationQueue))
			fmt.Printf("[miner] len(m.operationQueue) >= OP_THRESHOLD: %v\n", len(m.operationQueue) >= OP_THRESHOLD)
			
			if len(m.operationQueue) >= OP_THRESHOLD {
				difficulty = m.Settings.PoWDifficultyOpBlock
				for len(b.Operations) <= OP_THRESHOLD {
					fmt.Printf("[miner] Pre Dequeue len OperationQueue: %v\n", len(m.operationQueue))
					op := <- m.operationQueue
					fmt.Printf("[miner] Post Dequeue len OperationQueue: %v\n", len(m.operationQueue))
					validatedOp, err := m.ValidateOperation(op)
					fmt.Printf("[miner] validatedOp: %v\n", validatedOp)
					if err != nil {
						err = fmt.Errorf("unable to validate mining op: %v", err)
						fmt.Printf("[miner] validateOp error: %v\n", err)
						return err
					}
					if !validatedOp {
						continue
					}
					b.Operations = append(b.Operations, op)
				}
			} else {
				difficulty = m.Settings.PoWDifficultyNoOpBlock
			}

			fmt.Printf("[miner] Starting Mining: %v\n", b)
			err = b.Mine(difficulty)
			hash, err := b.Hash()
			fmt.Printf("[miner] Done Mining: %v with %s\n", b, hash)

			select {
			case <-m.earlyExitSignal:
				// readd the operations for the future
				for _, op := range b.Operations {
					// avoid a block. Toss op if full
					if len(m.operationQueue) < MAX_WAITING_OPS {
						m.operationQueue <- op
					}
				}
				minersGroup.Done()
				fmt.Printf("[miner] Early Exit\n")
				return nil
			default:
				_ = m.Blockchain.AppendBlock(b, m.Settings)
				err := m.DisseminateBlock(b)
				if err != nil {
					fmt.Printf("dissemination created error: %v", err)
					return err
				}
				m.OnNewBlock(*b)
			}
		}
	}
}

// this stops the process of mining blocks
// Commands the lower level threads to stop.
// Waitgroup finishes when these exit
func (m *Miner) StopMining() (err error) {
	fmt.Printf("[miner] Stopping Mining by command\n")
	close(m.earlyExitSignal)
	minersGroup.Wait()
	fmt.Printf("[miner] Stopped\n")
	return nil
}

func (m *Miner) TestEarlyExit() {
	time.Sleep(20000 * time.Millisecond)
	fmt.Printf("[miner] Killing...\n")
	err := m.StopMining()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (m *Miner) ValidateOperation(op *blockartlib.Operation) (bool, error) {
	// check sigs
	if op.ShapeHash != op.CalculateSig() {
		fmt.Printf("[miner#ValidateOperation] op '%v' fails sig check\n", op)
		return false, nil
	}
	// check drawable (implicitly already drawn)
	fmt.Printf("[miners#ValidateOperation] op: '%v'\n", op)
	fmt.Printf("[miners#ValidateOperation] op.ShapeSVGString: '%v'\n", op.ShapeSVGString)
	validOps, invalidOps, err := DrawOperations([]blockartlib.Operation{*op}, m.Settings.CanvasSettings)
	fmt.Printf("[miners#ValidateOperation] DrawOperations err: '%v'\n", err)
	if err != nil {
		return false, fmt.Errorf("unable to validate operation %v: %v", op, err)
	}
	if len(validOps) != 1 || len(invalidOps) != 0 || validOps[op.ShapeHash] != *op {
		fmt.Printf("[miner#ValidateOperation] op '%v' fails drawable check\n", op)
		fmt.Printf("[miner#ValidateOperation] validOps: '%v', invalidOps: '%v'\n", validOps, invalidOps)
		return false, nil
	}
	return true, nil
}

// validates incoming block from other miner
func (m *Miner) ValidNewBlock(b *Block) (valid bool, err error) {
	blockValid, err := b.Valid(m.Settings.PoWDifficultyOpBlock, m.Settings.PoWDifficultyNoOpBlock)
	if err != nil {
		return false, fmt.Errorf("Unable to validate block: %v", err)
	}
	if !blockValid {
		return false, nil
	}
	// check if block is in tree
	present, err := m.Blockchain.BlockPresent(b)
	if err != nil {
		return false, fmt.Errorf("Unable validate block: %v", err)
	}
	if present {
		return false, nil
	} else {
		// not present, is parent in tree
		parentBlock, err := m.Blockchain.FindBlockByHash(b.ParentHash)
		if err != nil {
			return false, fmt.Errorf("Unable validate block: %v", err)
		}
		if parentBlock != nil {
			// is found parent internally valid
			parentValid, err := parentBlock.Valid(m.Settings.PoWDifficultyOpBlock, m.Settings.PoWDifficultyNoOpBlock)
			if err != nil {
				return false, fmt.Errorf("Unable validate block: %v", err)
			}
			if parentValid {
				return true, nil
			}
		// failing that, check its ancestors in the forest, and those pass ValidBlock
		} else{
			forestParent :=  m.BlockForest[b.ParentHash] 
			if forestParent != nil {
				// internally consistentgi
				parentValid, err := forestParent.Valid(m.Settings.PoWDifficultyOpBlock, m.Settings.PoWDifficultyNoOpBlock)
				if err != nil {
					return false, fmt.Errorf("Unable validate block: %v", err)
				}
				if !parentValid {
					return false, nil
				}
				// forestParent validNewBlock too?
				forestParentValid, err := m.ValidNewBlock(m.BlockForest[b.ParentHash])
				if err != nil {
					return false, fmt.Errorf("Unable validate block: %v", err)
				}
				if forestParentValid {
					// TODO: Add forest Parent
					return true, nil
				}
			}
		}
		return false, nil
	}
}

/////// functions to interact with other miners

func (m *Miner) OpenNeighbourConnections() (err error) {
	/* Opens RPC connections to neighbouring Miners and fills in the
	   RPCClient field in the corresponding MinerConnection struct */
	for i, neighbour := range m.Neighbours {
		if neighbour.Alive {
			continue
		}
		fmt.Println("Before open RPC to neighbour: ", neighbour.Addr.String())
		neighbour.RPCClient, err = rpc.Dial("tcp", neighbour.Addr.String())
		if err != nil {
			deleteNeighbour(m, i)
			return nil
		}
		fmt.Printf("Opened RPC connection to neighbour with tcpAddr %s\n", neighbour.Addr.String())
	}

	return nil
}

// TODO: Actually handle the case where the blockchain we choose is invalid
func (m *Miner) ConnectToNeighbourMiners(localAddr *net.TCPAddr) (bestNeighbour net.TCPAddr, err error) {
	/* Makes the RPC call to register itself to neighbouring miners.
	   Neighbours will respond with the length of their Blockchain;
	   does NOT currently account for the fact that the given chain
	   may be invalid.

	   RETURNS: the net.TCPAddress of the neighbour with the longest
	   blockchain depth
	*/

	// Connect to each neighbour miner and keep track of the one with the largest depth
	var bestMinerAddr net.TCPAddr
	largestDepth := m.Blockchain.BC.LastNode.Current.Depth
	depth := 0

	fmt.Println("Our address before sending RPC call: ", localAddr.String())
	for i, connection := range m.Neighbours {
		/*fmt.Println("DISCONNECT!!!")
		time.Sleep(4*time.Second)*/
		err = connection.RPCClient.Call("MinerInstance.ConnectNewNeighbour", localAddr, &depth)
		if err != nil {
			// TODO: Should we just ignore this miner then and move on to the next one?
			erro := deleteNeighbour(m, i)
			blockartlib.CheckErr(erro)
			//return net.TCPAddr{}, err
			continue
		}

		//if (depth >= largestDepth) {
		if (depth > largestDepth) {
			largestDepth = depth
			bestMinerAddr = connection.Addr
		}
	}

	return bestMinerAddr, nil
}

// TODO: Actually handle the case where the blockchain we choose is invalid
/*func (m *Miner) RequestBCStorageFromNeighbour(neighbourAddr *net.TCPAddr) (err error) {
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	treeArray := make([][]byte, 0)
	for i, v := range m.Neighbours {
		if v.Addr.String() == neighbourAddr.String() {
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
}*/

func (m *Miner) RequestBCStorageFromNeighbour(neighbourAddr *net.TCPAddr, queue *[]string) (err error) {
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	//treeArray := make([][]byte, 0)
	for i, v := range m.Neighbours {
		if v.Addr.String() == neighbourAddr.String() {
			err := v.RPCClient.Call("MinerInstance.DisseminateTree", true, &queue)
			if err != nil {
				deleteNeighbour(m, i)
			}
		}
	}
	if len(*queue) != 0 {
		reconstructTree(m, neighbourAddr, queue)
	}
	return nil
}

// sends out the block to other miners
func (m *Miner) DisseminateBlock(block *Block) (err error) {
	// TODO: notify waiting artNodes if your block is op number of nodes deep now
	// TODO: Not sure this is the right spot for this.
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	for _,v := range m.Neighbours {
		marshalledBlock, err := block.MarshallBinary()
		fmt.Println("Marshalled block in disseminateBlock: ", marshalledBlock)
		blockartlib.CheckErr(err)
		var b bool
		err = v.RPCClient.Call("MinerInstance.ReceiveBlockFromNeighbour", &marshalledBlock, &b)
		if !b {
			fmt.Println("Bad block") // TODO: think what to do in this case
		}
	}
	return err
}

func (m *Miner) AddDisseminatedBlock(b *Block) {
	valid, err := m.ValidNewBlock(b)
	if err != nil {
		fmt.Printf("Error in AddDisseminatedBlock: %v", err)
		return
	}
	if valid {
		// Add to blockchain
		treeSwitch := m.Blockchain.AppendBlock(b, m.Settings)
		if treeSwitch {
			// blocks until complete
			m.StopMining()
			// mines on the new longest chain
			m.StartMining()
		}
		m.OnNewBlock(*b)
	}
}

func (m *Miner) DisseminateOperation(op Operation) (err error) {
	// TODO: If any changes are made in disseminate block, repeat here
	for _, v := range m.Neighbours {
		marshalledOp, err := op.Marshall()
		blockartlib.CheckErr(err)
		var b bool
		err = v.RPCClient.Call("MinerInstance.ReceiveOperationFromNeighbour", &marshalledOp, &b)
		if !b {
			fmt.Println("bad op")
		}
	}
	return err
}

func (m *Miner) GetShapeHash(op *blockartlib.Operation) (shapeHash string, err error) {
	fmt.Println("MINERLIB: Running GetShapeHash.")
	fmt.Println("MINERLIB: Trying to find ArtNodeConnection using the ArtNode's public key")
	artNodeConn, err := m.FindArtNodeConnection(op.ArtNodePubKey)
	if err != nil {
		return "", fmt.Errorf("unable to get shape hash: %v", err)
	}
	if artNodeConn == nil {
		return "", fmt.Errorf("unable to locate the art node")
	}
	fmt.Println("MINERLIB: Found ArtNode connection; returning shape hash response")
	// TODO: Uncomment this when channels work
	// blocks until a value comes into ShapeHashResponse
	// return <- artNodeConn.ShapeHashResponse, nil
	return "thisIsAHash", nil
}

func (m *Miner) OnNewBlock(b Block) {
	for _, op := range b.Operations {
		// iterate through slice 0
		for _, doneOp := range m.OpValidateList[0] {
			artNodeConn, err := m.FindArtNodeConnection(doneOp.ArtNodePubKey)
			if err != nil {
				fmt.Printf("ERROR in OnNewBlock: %v", err)
				return
			}
			if artNodeConn == nil {
				continue
			}
			// TODO: Check if returning the right portion of op
			// fill channel for associated art node
			artNodeConn.ShapeHashResponse <- doneOp.ShapeHash
		}
		if m.HasArtNode(op.ArtNodePubKey) {
			m.OpValidateList[op.ValidateBlockNum] = append(m.OpValidateList[op.ValidateBlockNum], op)
		}
	}
	if len(m.OpValidateList) >= 1 {
		m.OpValidateList = m.OpValidateList[1:]
	}
}

func (m *Miner) FindArtNodeConnection(artNodePublicKey string) (anc *ArtNodeConnection, err error) {
	fmt.Println("MINER: Running FindArtNodeConnection with public key", artNodePublicKey)
	for i, an := range m.ArtNodes {
		fmt.Printf("Index %d, key %s\n", i, an.ArtNodePubKey)
		if an.ArtNodePubKey == artNodePublicKey {
			return an, nil
		}
	}
	return nil, nil
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

/*func (m *Miner) MarshallTree (result *[][]byte, node *BCTreeNode) *[][]byte{
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	tree := m.Blockchain.BCT
	//genBlock := tree.GenesisNode.BlockResiding
	//marshalledGenBlock, _ := genBlock.MarshallBinary()
	//*result = append(*result, marshalledGenBlock)
	//for len(tree.GenesisNode.Children) != 0 {
	if node == nil {
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
				*result = append(*result, marshalledBlock[0])
				b:= m.MarshallTree(result, node)
				*result = append(*result, *b...)
			}
		}
	}

	return result
}*/

func (m *Miner) MarshallTree (result *[]string, node *BCTreeNode) []string{
	gob.Register(&Block{})
	gob.Register(elliptic.P384())
	tree := m.Blockchain.BCT
	if node == nil {
		for range tree.GenesisNode.Children {
			children, err := m.Blockchain.GetChildrenNodes(tree.GenesisNode.CurrentHash)
			blockartlib.CheckErr(err)
			for _, v := range children {
				node := FindBCTreeNode(tree.GenesisNode,v)
				block := node.BlockResiding
				blockHash, err := block.Hash()
				if err != nil {
					fmt.Println("Hashing on tree request unsuccessful")
					return *result
				}
				*result = append(*result, blockHash)
				b:= m.MarshallTree(result, node)
				*result = append(*result, b...)
			}
		}
	}

	return *result
}

func (m *Miner) AddMinersToList(lom *[]net.Addr) (err error) {
	if len(*lom) == 0 {
		return nil
	} else if len(m.Neighbours) == 0 {
		for _, val := range *lom {
			addMinerToList(m, val)
		}
	} else if len(m.Neighbours) > 0 {
		for _, val := range *lom {
			if len(m.Neighbours) == 256 {
				return nil
			}
			if !isMinerInList(m, val) {
				addMinerToList(m, val)
			}
		}
	}
	return nil
}

func (m *Miner) HasArtNode(artNodePubKey string) bool {
	hasArtNode := false
	for _, conn := range m.ArtNodes {
		if conn.ArtNodePubKey == artNodePubKey {
			hasArtNode = true
			break
		}
	}
	return hasArtNode
}

func addMinerToList(m *Miner, addr net.Addr) error {
	var newNeighbour = MinerConnection{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return err
	}
	newNeighbour.Addr = *tcpAddr
	m.Neighbours = append(m.Neighbours, &newNeighbour)
	return nil
}

func isMinerInList(m *Miner, addr net.Addr) bool {
	for _, v := range m.Neighbours {
		if v.Addr.String() == addr.String() {
			return true
		}
	}
	return false
}

func deleteNeighbour (m *Miner, index int) error {
	buf := m.Neighbours[:index]
	m.Neighbours = append(buf, m.Neighbours[index+1:]...)
	return nil
}

/*func reconstructTree(m *Miner, tree *[][]byte) {
	t := *tree
	genBlock := m.CreateGenesisBlock()
	fmt.Println("tree: ", t)
	m.Blockchain = NewBlockchainStorage(genBlock, m.Settings)
	fmt.Println("New Blockchain: ", m.Blockchain.BCT.GenesisNode.CurrentHash)
	t = t[1:]
	for _,v := range t {
		var temp [][]byte
		temp[0] = v
		b, err := UnmarshallBinary(temp)
		fmt.Println("the block received: ", b)
		if err != nil {
			fmt.Println("unmarshalling failed")
			return
		}
		valid, err := m.ValidNewBlock(b)
		blockartlib.CheckErr(err)
		if err != nil || !valid{
			fmt.Printf("Invalid block: %v", err)
			return
		}
		m.Blockchain.AppendBlock(b, m.Settings)
	}
}*/

func reconstructTree(m *Miner, senderAddr *net.TCPAddr, queue *[]string) {
	q := *queue
	genBlock := m.CreateGenesisBlock()
	fmt.Println("queue: ", q)
	m.Blockchain = NewBlockchainStorage(genBlock, m.Settings)
	fmt.Println("New Blockchain: ", m.Blockchain.BCT.GenesisNode.CurrentHash)
	q = q[1:]
	blockArr := make([][]byte, 0)
	var caller *rpc.Client
	for _,v := range m.Neighbours {
		if v.Addr.String() == senderAddr.String() {
			caller = v.RPCClient
			break
		}
	}

	for _,v := range q {
		err := caller.Call("MinerInstance.GiveBlock", &v, &blockArr)
		fmt.Println("the block received: ", blockArr)
		b, err := UnmarshallBinary(blockArr)
		if err != nil {
			fmt.Println("unmarshalling failed")
			return
		}
		valid, err := m.ValidNewBlock(b)
		blockartlib.CheckErr(err)
		if err != nil || !valid{
			fmt.Printf("Invalid block: %v", err)
			return
		}
		m.Blockchain.AppendBlock(b, m.Settings)
	}
}
