package miner

type Miner struct {
	HBeatAddr string
	Nbrs      [256]*Miner
	
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