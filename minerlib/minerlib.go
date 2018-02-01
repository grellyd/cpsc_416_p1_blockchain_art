package minerlib

type Miner struct {
	HBeatAddr string
	Nbrs      [256]int
	MAddr     string
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
