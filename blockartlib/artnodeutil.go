package blockartlib

type AN struct {
	MinerID		int //keep reference to the connected miner
	MinerAddr 	string
	PrivKey 	string
	PubKey 		string
	MinerAlive 	bool
}

func (an *AN) Connect(miner string) (err error) {
	return nil
}

func (an *AN) MakeDrawRequest() (err error) {
	return err
}