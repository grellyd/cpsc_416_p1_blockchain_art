package blockartlib

/*
Artnode that communicates with the client app and the miner
*/

type AN struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
}


// CANVAS INTERFACE FUNCTIONS
func (an *AN) AddShape(validateNum uint8, shapeType ShapeType, shapeSvgString string, fill string, stroke string) (shapeHash string, blockHash string, inkRemaining uint32, err error) {
	return shapeHash, blockHash, inkRemaining, err
}

func (an *AN) GetSvgString(shapeHash string) (svgString string, err error) {
	return svgString, err
}

func (an *AN) GetInk() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *AN) DeleteShape(validateNum uint8, shapeHash string) (inkRemaining uint32, err error) {
	return inkRemaining, err
}

func (an *AN) GetShapes(blockHash string) (shapeHashes []string, err error) {
	return shapeHashes, err
}

func (an *AN) GetGenesisBlock() (blockHash string, err error) {
	return blockHash, err
}

func (an *AN) GetChildren(blockHash string) (blockHashes []string, err error) {
	return blockHashes, err
}
func (an *AN) CloseCanvas() (inkRemaining uint32, err error) {
	return inkRemaining, err
}

// MINER INTERACTION FUNCTIONS
func (an *AN) Connect(minerAddr, pubKey, privKey string) (err error) {
	
	return nil
}

func (an *AN) MakeDrawRequest() (err error) {
	return err
}
