package minerlib

type Miner struct {
	HBeatAddr string
	Nbrs      [256]int
	ThrNbrs   int // min threshold below which we request new neighbours
	//MAddr     string
	ServAddr  string
	ServAlive bool
	//PubKey    string
	//PrivKey   string
	//InkLevel  int
	//ANs       []int // maybe []*AN
}

type Block struct {
}

func (m *Miner) ValidateNewArtlib() (err error) {
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

/////// functions to interact with server

// retrieves settings from server
func (m *Miner) RetrieveSettings() (err error) {
	return nil
}

// requests another miner's ID (info) from the server
func (m *Miner) RequestMiner() (err error) {
	return nil
}

func (m *Miner) SendHeartbeat() (err error) {
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
