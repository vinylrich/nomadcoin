package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
	Height   int
}

func (b *Block) fromBytes(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleError(decoder.Decode(b))
}
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)

	block.Hash = fmt.Sprintf("%s", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}

	block := &Block{}
	block.fromBytes(blockBytes)

	return block, nil
}
