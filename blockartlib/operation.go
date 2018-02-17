package blockartlib

import (
	"encoding/gob"
	"fmt"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/binary"
)

type OperationType uint32

const (
	NOP OperationType = iota
	DRAW
	DELETE
)

type Operation struct {
	Type OperationType
	OperationNumber uint32
	Shape ShapeType
	Fill string // Can be "transparent" or "filled"
	Stroke string
	ShapeSVGString string
	ArtNodePubKey string
	ValidateBlockNum uint8
	ShapeHash string
	Nonce uint32
}

// Let the OpSig be the MD5 Hash of the operation type, operation number, and the ArtNode's Public Key
func (o *Operation) GenerateSig() error {
	o.ShapeHash = o.CalculateSig()
	return nil
}

func (o *Operation) CalculateSig() (string) {
	data := []byte{}
	h := md5.New()
	data = append(data, uint32AsByteArr(uint32(o.Type))...)
	data = append(data, uint32AsByteArr(o.OperationNumber)...)
	h.Write(append(data, []byte(o.ArtNodePubKey)...))
	return hex.EncodeToString(h.Sum(nil))
}

func (o *Operation)Marshall() (data []byte, err error) {
	gob.Register(&Operation{})
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(*o)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall operation: %v", err)
	}
	return buff.Bytes(), nil
}

func OperationUnmarshall(data []byte) (o *Operation, err error) {
	// o = &Operation{}
	gob.Register(o)
	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(o)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall operation: %v", err)
	}
	return o, nil
}

// duplicated code from minerlib/compute
// TODO: extract out to marshall_utils
func uint32AsByteArr(val uint32) []byte {
	a := make([]byte, 4)
    binary.LittleEndian.PutUint32(a, val)
	return a
}
